package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"terraform-provider-verity/internal/auth"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &verityProvider{}

type verityProvider struct {
	version string
}

type providerContext struct {
	client       *openapi.APIClient
	tokenManager *auth.TokenManager
	config       *openapi.Configuration
	credentials  struct {
		username string
		password string
	}
	responseCache  map[string]interface{}
	cacheMutex     sync.Mutex
	bulkOpsMgr     *utils.BulkOperationManager
	tickChannel    chan struct{}
	debounceTimer  *time.Timer
	debounceActive bool
	debounceMutex  sync.Mutex
}

type verityProviderModel struct {
	URI      types.String `tfsdk:"uri"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &verityProvider{
			version: version,
		}
	}
}

func (p *verityProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "verity"
}

func (p *verityProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Verity API",
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				Description: "The base URL of the API",
				Optional:    true,
				Sensitive:   true,
			},
			"username": schema.StringAttribute{
				Description: "API username",
				Optional:    true,
				Sensitive:   true,
			},
			"password": schema.StringAttribute{
				Description: "API password",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *verityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config verityProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uri := config.URI.ValueString()
	if uri == "" {
		uri = os.Getenv("TF_VAR_uri")
		tflog.Debug(ctx, "URI not provided in configuration, using environment variable")
	}

	username := config.Username.ValueString()
	if username == "" {
		username = os.Getenv("TF_VAR_username")
		tflog.Debug(ctx, "Username not provided in configuration, using environment variable")
	}

	password := config.Password.ValueString()
	if password == "" {
		password = os.Getenv("TF_VAR_password")
		tflog.Debug(ctx, "Password not provided in configuration, using environment variable")
	}

	if uri == "" {
		resp.Diagnostics.AddError(
			"Missing API URI",
			"The provider cannot create the Verity API client as the URI is missing. "+
				"Set the uri attribute in the provider configuration or "+
				"set the TF_VAR_uri environment variable.",
		)
		return
	}

	if username == "" {
		resp.Diagnostics.AddError(
			"Missing API Username",
			"The provider cannot create the Verity API client as the username is missing. "+
				"Set the username attribute in the provider configuration or "+
				"set the TF_VAR_username environment variable.",
		)
		return
	}

	if password == "" {
		resp.Diagnostics.AddError(
			"Missing API Password",
			"The provider cannot create the Verity API client as the password is missing. "+
				"Set the password attribute in the provider configuration or "+
				"set the TF_VAR_password environment variable.",
		)
		return
	}

	apiConfig := openapi.NewConfiguration()

	jar, err := cookiejar.New(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create cookie jar",
			fmt.Sprintf("Failed to create cookie jar: %v", err),
		)
		return
	}

	tokenManager := auth.NewTokenManager(jar)

	apiConfig.HTTPClient = &http.Client{
		Jar: jar,
	}

	baseURL := uri
	tflog.Debug(ctx, "Configuring provider", map[string]interface{}{
		"base_url": baseURL,
	})

	baseURL = strings.TrimRight(baseURL, "/")

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid URI",
			fmt.Sprintf("Invalid URI: %v", err),
		)
		return
	}

	apiConfig.Host = parsedURL.Host
	apiConfig.Scheme = parsedURL.Scheme
	tflog.Debug(ctx, "Parsed URL configuration", map[string]interface{}{
		"host":   apiConfig.Host,
		"scheme": apiConfig.Scheme,
		"url":    parsedURL.String(),
	})

	if apiConfig.Scheme == "" {
		resp.Diagnostics.AddError(
			"Invalid URI Scheme",
			"URI must start with http:// or https://",
		)
		return
	}

	serverURL := fmt.Sprintf("%s://%s/api", apiConfig.Scheme, apiConfig.Host)
	apiConfig.Servers = openapi.ServerConfigurations{
		{
			URL: serverURL,
		},
	}
	tflog.Debug(ctx, "Server configuration", map[string]interface{}{
		"url": serverURL,
	})

	apiConfig.Debug = true
	client := openapi.NewAPIClient(apiConfig)

	provCtx := &providerContext{
		client:         client,
		tokenManager:   tokenManager,
		config:         apiConfig,
		responseCache:  make(map[string]interface{}),
		debounceActive: true,
	}
	provCtx.credentials.username = username
	provCtx.credentials.password = password

	bulkManager := utils.GetBulkOperationManager(client, clearCache, provCtx)
	tflog.Info(ctx, "Initialized bulk operation manager with manual batching mode")

	provCtx.bulkOpsMgr = bulkManager

	provCtx.initBulkOpsTicker(ctx)

	if err := authenticate(ctx, provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Authentication Failed",
			fmt.Sprintf("Failed to authenticate with Verity API: %v", err),
		)
		return
	}

	resp.ResourceData = provCtx
	resp.DataSourceData = provCtx
}

func (p *verityProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVerityTenantResource,
		NewVerityGatewayResource,
		NewVerityServiceResource,
		NewVerityEthPortProfileResource,
		NewVerityEthPortSettingsResource,
		NewVerityBundleResource,
		NewVerityLagResource,
		NewVerityGatewayProfileResource,
		NewVerityOperationStageResource,
	}
}

func (p *verityProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVerityStateImporterDataSource,
	}
}

func (p *providerContext) initBulkOpsTicker(ctx context.Context) {
	p.tickChannel = make(chan struct{}, 100)

	go func() {
		for range p.tickChannel {
			p.debounceMutex.Lock()
			if p.debounceTimer != nil {
				p.debounceTimer.Stop()
			}

			// when no new ticks arrive for 3 seconds, execute operations
			p.debounceTimer = time.AfterFunc(3*time.Second, func() {
				tflog.Debug(ctx, "Bulk operation debounce timer expired, executing pending operations")
				if diags := p.bulkOpsMgr.ExecuteAllPendingOperations(ctx); diags != nil {
					tflog.Error(ctx, "Failed to execute pending bulk operations", map[string]interface{}{
						"error": diags,
					})
				}
			})
			p.debounceMutex.Unlock()
		}
	}()

	tflog.Debug(ctx, "Bulk operations ticker initialized")
}

func (p *providerContext) NotifyOperationAdded() {
	if p.debounceActive {
		p.tickChannel <- struct{}{}
		tflog.Debug(context.Background(), "Operation added, tick sent to debounce system")
	}
}

func authenticate(ctx context.Context, provCtx *providerContext) error {
	token, needsRefresh := provCtx.tokenManager.GetToken()
	if !needsRefresh {
		u, _ := url.Parse(provCtx.config.Servers[0].URL)
		provCtx.config.HTTPClient.Jar.SetCookies(u, []*http.Cookie{
			{
				Name:  "ivn_api",
				Value: token,
			},
		})
		return nil
	}

	auth := openapi.NewAuthPostRequestAuth(
		provCtx.credentials.username,
		provCtx.credentials.password,
	)
	authReq := openapi.NewAuthPostRequest()
	authReq.SetAuth(*auth)

	resp, err := provCtx.client.AuthorizationAPI.AuthPost(ctx).AuthPostRequest(*authReq).Execute()
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Token == "" {
		return fmt.Errorf("no token found in response")
	}

	provCtx.tokenManager.SetToken(result.Token, 24*time.Hour)

	u, _ := url.Parse(provCtx.config.Servers[0].URL)
	provCtx.config.HTTPClient.Jar.SetCookies(u, []*http.Cookie{
		{
			Name:  "ivn_api",
			Value: result.Token,
		},
	})

	return nil
}

func ensureAuthenticated(ctx context.Context, m interface{}) error {
	provCtx := m.(*providerContext)
	return authenticate(ctx, provCtx)
}

func getCachedResponse(ctx context.Context, m interface{}, cacheKey string, apiCall func() (interface{}, error), forceRefresh ...bool) (interface{}, error) {
	provCtx := m.(*providerContext)

	provCtx.cacheMutex.Lock()
	defer provCtx.cacheMutex.Unlock()

	shouldRefresh := false
	if len(forceRefresh) > 0 && forceRefresh[0] {
		shouldRefresh = true
		tflog.Debug(ctx, "Force refreshing cache for "+cacheKey)
	}

	if !shouldRefresh {
		if cachedResp, ok := provCtx.responseCache[cacheKey]; ok {
			tflog.Debug(ctx, "Using cached response for "+cacheKey)
			return cachedResp, nil
		}
	}

	resp, err := apiCall()
	if err != nil {
		return nil, err
	}

	provCtx.responseCache[cacheKey] = resp
	tflog.Debug(ctx, "Cached new response for "+cacheKey)

	return resp, nil
}

func clearCache(ctx context.Context, m interface{}, cacheKey string) {
	provCtx := m.(*providerContext)

	provCtx.cacheMutex.Lock()
	defer provCtx.cacheMutex.Unlock()
	if cacheKey == "" {
		provCtx.responseCache = make(map[string]interface{})
		tflog.Debug(ctx, "Cleared entire response cache")
	} else if _, ok := provCtx.responseCache[cacheKey]; ok {
		delete(provCtx.responseCache, cacheKey)
		tflog.Debug(ctx, "Cleared cache for "+cacheKey)
	}
}
