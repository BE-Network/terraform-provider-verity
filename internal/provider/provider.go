package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	mode           string
	apiVersion     string
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
	Mode     types.String `tfsdk:"mode"`
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
			"mode": schema.StringAttribute{
				Description: "Mode to operate in: 'datacenter' (default) or 'campus'",
				Optional:    true,
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

	mode := config.Mode.ValueString()
	if mode == "" {
		mode = os.Getenv("TF_VAR_mode")
		tflog.Debug(ctx, "Mode not provided in configuration, checking environment variable")
	}

	// Default to 'datacenter' mode if not specified
	if mode == "" {
		mode = "datacenter"
		tflog.Debug(ctx, "Mode not provided in configuration or environment, defaulting to 'datacenter'")
	} else if mode != "datacenter" && mode != "campus" {
		resp.Diagnostics.AddError(
			"Invalid Mode",
			"The mode must be either 'datacenter' or 'campus'. "+
				"Got: "+mode,
		)
		return
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
		config:         apiConfig,
		client:         client,
		tokenManager:   tokenManager,
		responseCache:  make(map[string]interface{}),
		mode:           mode,
		debounceActive: true,
	}

	provCtx.credentials.username = username
	provCtx.credentials.password = password

	tflog.Info(ctx, "Configuring provider with mode: "+mode)

	bulkManager := utils.GetBulkOperationManager(client, clearCache, provCtx, mode)
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

	tflog.SetField(ctx, "verity_mode", provCtx.mode)

	apiVersion, err := getApiVersion(ctx, provCtx)
	if err != nil {
		if strings.Contains(err.Error(), "mode mismatch") {
			resp.Diagnostics.AddError(
				"Provider Mode Configuration Error",
				err.Error(),
			)
			return
		}

		tflog.Warn(ctx, "Failed to get API version, using default compatibility", map[string]interface{}{
			"error": err.Error(),
		})
		// Set a default version if we couldn't get the actual version
		apiVersion = "6.4.0"
	}

	provCtx.apiVersion = apiVersion

	ctxWithProviderData := context.WithValue(ctx, "providerData", provCtx)

	resp.DataSourceData = provCtx
	resp.ResourceData = provCtx

	tflog.Info(ctxWithProviderData, "Provider configured", map[string]interface{}{
		"mode":        provCtx.mode,
		"api_version": provCtx.apiVersion,
	})
}

func getApiVersion(ctx context.Context, provCtx *providerContext) (string, error) {
	const defaultVersion = "6.4"

	if err := authenticate(ctx, provCtx); err != nil {
		return "", fmt.Errorf("authentication failed when getting API version: %w", err)
	}

	tflog.Debug(ctx, "Fetching API version via OpenAPI client")

	httpResp, err := provCtx.client.VersionAPI.VersionGet(ctx).Execute()

	if err != nil {
		tflog.Error(ctx, "Failed to execute API version request via OpenAPI client, defaulting.", map[string]interface{}{
			"error":           err,
			"default_version": defaultVersion,
		})
		return defaultVersion, nil // Return default version, no error to the caller of getApiVersion itself.
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		tflog.Error(ctx, "Failed to read API version response body (OpenAPI client), defaulting.", map[string]interface{}{
			"error":           err,
			"default_version": defaultVersion,
		})
		return defaultVersion, nil
	}

	if httpResp.StatusCode != http.StatusOK {
		tflog.Error(ctx, "API version request via OpenAPI client returned non-OK status, defaulting.", map[string]interface{}{
			"status_code":     httpResp.StatusCode,
			"response_body":   string(bodyBytes),
			"default_version": defaultVersion,
		})
		return defaultVersion, nil
	}

	var versionPayload struct {
		Version    string `json:"version"`
		Datacenter *bool  `json:"datacenter,omitempty"`
	}
	if err := json.Unmarshal(bodyBytes, &versionPayload); err != nil {
		tflog.Error(ctx, "Failed to unmarshal API version JSON response (OpenAPI client), defaulting.", map[string]interface{}{
			"error":           err,
			"response_body":   string(bodyBytes),
			"default_version": defaultVersion,
		})
		return defaultVersion, nil
	}

	if versionPayload.Version == "" {
		tflog.Warn(ctx, "API version string is empty in response (OpenAPI client), defaulting.", map[string]interface{}{
			"response_body":   string(bodyBytes),
			"default_version": defaultVersion,
		})
		return defaultVersion, nil
	}

	if versionPayload.Datacenter != nil {
		systemIsDatacenter := *versionPayload.Datacenter
		configuredIsDatacenter := provCtx.mode == "datacenter"

		if systemIsDatacenter != configuredIsDatacenter {
			var systemMode, configuredMode string
			if systemIsDatacenter {
				systemMode = "datacenter"
			} else {
				systemMode = "campus"
			}
			configuredMode = provCtx.mode

			return "", fmt.Errorf("mode mismatch: provider is configured for '%s' mode but the system is running in '%s' mode. Please update the provider configuration to match the actual system type", configuredMode, systemMode)
		}

		tflog.Info(ctx, "Mode validation successful", map[string]interface{}{
			"configured_mode": provCtx.mode,
			"system_mode":     systemIsDatacenter,
		})
	} else {
		tflog.Debug(ctx, "No datacenter field in version response, skipping mode validation")
	}

	tflog.Info(ctx, "Successfully fetched API version via OpenAPI client", map[string]interface{}{
		"version": versionPayload.Version,
	})
	return versionPayload.Version, nil
}

func (p *verityProvider) Resources(ctx context.Context) []func() resource.Resource {
	providerData, ok := ctx.Value("providerData").(*providerContext)
	if !ok {
		tflog.Warn(ctx, "Provider context not available, returning all resources")
		return getAllResources()
	}

	allResources := getAllResources()
	compatibleResources := utils.FilterResourcesByMode(ctx, allResources, providerData.mode, providerData.apiVersion)

	return compatibleResources
}

func getAllResources() []func() resource.Resource {
	return []func() resource.Resource{
		NewVerityOperationStageResource,
		NewVerityTenantResource,
		NewVerityGatewayResource,
		NewVerityServiceResource,
		NewVerityEthPortProfileResource,
		NewVerityEthPortSettingsResource,
		NewVerityBundleResource,
		NewVerityLagResource,
		NewVerityGatewayProfileResource,
		NewVerityACLV4Resource,
		NewVerityACLV6Resource,
		NewVerityBadgeResource,
		NewVerityAuthenticatedEthPortResource,
		NewVerityDeviceVoiceSettingsResource,
		NewVerityPacketBrokerResource,
		NewVerityPacketQueueResource,
		NewVerityServicePortProfileResource,
		NewVerityVoicePortProfileResource,
		NewVeritySwitchpointResource,
		NewVerityDeviceControllerResource,
		NewVerityAsPathAccessListResource,
		NewVerityCommunityListResource,
		NewVerityDeviceSettingsResource,
		NewVerityExtendedCommunityListResource,
		NewVerityIpv4ListResource,
		NewVerityIpv4PrefixListResource,
		NewVerityIpv6ListResource,
		NewVerityIpv6PrefixListResource,
		NewVerityRouteMapClauseResource,
		NewVerityRouteMapResource,
		NewVeritySfpBreakoutResource,
		NewVeritySiteResource,
		NewVerityPodResource,
		NewVerityPortAclResource,
		NewVeritySflowCollectorResource,
		NewVerityDiagnosticsProfileResource,
		NewVerityDiagnosticsPortProfileResource,
		NewVerityPBRoutingResource,
		NewVerityPBRoutingACLResource,
		NewVeritySpinePlaneResource,
		NewVerityGroupingRuleResource,
		NewVerityThresholdGroupResource,
		NewVerityThresholdResource,
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
