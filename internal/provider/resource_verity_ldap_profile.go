package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityLdapProfileResource{}
	_ resource.ResourceWithConfigure   = &verityLdapProfileResource{}
	_ resource.ResourceWithImportState = &verityLdapProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityLdapProfileResource{}
)

const ldapProfileResourceType = "ldapprofiles"
const ldapProfileTerraformType = "verity_ldap_profile"

func NewVerityLdapProfileResource() resource.Resource {
	return &verityLdapProfileResource{}
}

type verityLdapProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityLdapProfileResourceModel struct {
	Name                     types.String                         `tfsdk:"name"`
	Enable                   types.Bool                           `tfsdk:"enable"`
	BaseDn                   types.String                         `tfsdk:"base_dn"`
	BindDn                   types.String                         `tfsdk:"bind_dn"`
	BindPassword             types.String                         `tfsdk:"bind_password"`
	EncryptedBindPassword    types.String                         `tfsdk:"encrypted_bind_password"`
	LdapVersion              types.String                         `tfsdk:"ldap_version"`
	SslTlsMode               types.String                         `tfsdk:"ssl_tls_mode"`
	DefaultPort              types.Int64                          `tfsdk:"default_port"`
	SearchTimeLimit          types.Int64                          `tfsdk:"search_time_limit"`
	BindTimeLimit            types.Int64                          `tfsdk:"bind_time_limit"`
	IdleTimeLimit            types.Int64                          `tfsdk:"idle_time_limit"`
	RetransmitAttempts       types.Int64                          `tfsdk:"retransmit_attempts"`
	SearchScope              types.String                         `tfsdk:"search_scope"`
	NssBasePasswd            types.String                         `tfsdk:"nss_base_passwd"`
	NssBaseGroup             types.String                         `tfsdk:"nss_base_group"`
	NssBaseShadow            types.String                         `tfsdk:"nss_base_shadow"`
	NssBaseNetgroup          types.String                         `tfsdk:"nss_base_netgroup"`
	NssBaseSudoers           types.String                         `tfsdk:"nss_base_sudoers"`
	NssInitgroupsIgnoreUsers types.String                         `tfsdk:"nss_initgroups_ignore_users"`
	NssSkipMembers           types.Bool                           `tfsdk:"nss_skip_members"`
	PamFilter                types.String                         `tfsdk:"pam_filter"`
	PamLoginAttribute        types.String                         `tfsdk:"pam_login_attribute"`
	PamGroupDn               types.String                         `tfsdk:"pam_group_dn"`
	PamMemberAttribute       types.String                         `tfsdk:"pam_member_attribute"`
	SudoersBase              types.String                         `tfsdk:"sudoers_base"`
	SudoersSearchFilter      types.String                         `tfsdk:"sudoers_search_filter"`
	LdapServers              []verityLdapProfileLdapServerModel   `tfsdk:"ldap_servers"`
	AttributeMaps            []verityLdapProfileAttributeMapModel `tfsdk:"attribute_maps"`
}

type verityLdapProfileLdapServerModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	Server             types.String `tfsdk:"server"`
	Port               types.Int64  `tfsdk:"port"`
	UseType            types.String `tfsdk:"use_type"`
	Priority           types.Int64  `tfsdk:"priority"`
	SslTlsMode         types.String `tfsdk:"ssl_tls_mode"`
	RetransmitAttempts types.Int64  `tfsdk:"retransmit_attempts"`
	Index              types.Int64  `tfsdk:"index"`
}

func (m verityLdapProfileLdapServerModel) GetIndex() types.Int64 {
	return m.Index
}

type verityLdapProfileAttributeMapModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	MapName types.String `tfsdk:"map_name"`
	From    types.String `tfsdk:"from"`
	To      types.String `tfsdk:"to"`
	Index   types.Int64  `tfsdk:"index"`
}

func (m verityLdapProfileAttributeMapModel) GetIndex() types.Int64 {
	return m.Index
}

func (r *verityLdapProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_profile"
}

func (r *verityLdapProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

func (r *verityLdapProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity LDAP Profile.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable object.",
				Optional:    true,
				Computed:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "Base Distinguished Name to use for LDAP searches",
				Optional:    true,
				Computed:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "Distinguished Name with which to bind to the LDAP server",
				Optional:    true,
				Computed:    true,
			},
			"bind_password": schema.StringAttribute{
				Description: "Credentials with which to bind to the LDAP server",
				Optional:    true,
				Computed:    true,
			},
			"encrypted_bind_password": schema.StringAttribute{
				Description: "System-generated encrypted version of Bind Password",
				Optional:    true,
				Computed:    true,
			},
			"ldap_version": schema.StringAttribute{
				Description: "LDAP protocol version",
				Optional:    true,
				Computed:    true,
			},
			"ssl_tls_mode": schema.StringAttribute{
				Description: "Global TLS mode for LDAP connections",
				Optional:    true,
				Computed:    true,
			},
			"default_port": schema.Int64Attribute{
				Description: "Default LDAP server port",
				Optional:    true,
				Computed:    true,
			},
			"search_time_limit": schema.Int64Attribute{
				Description: "Search time limit, in seconds",
				Optional:    true,
				Computed:    true,
			},
			"bind_time_limit": schema.Int64Attribute{
				Description: "Bind/connect time limit, in seconds",
				Optional:    true,
				Computed:    true,
			},
			"idle_time_limit": schema.Int64Attribute{
				Description: "NSS idle connection time limit, in seconds",
				Optional:    true,
				Computed:    true,
			},
			"retransmit_attempts": schema.Int64Attribute{
				Description: "Number of retransmit attempts",
				Optional:    true,
				Computed:    true,
			},
			"search_scope": schema.StringAttribute{
				Description: "Default LDAP search scope",
				Optional:    true,
				Computed:    true,
			},
			"nss_base_passwd": schema.StringAttribute{
				Description: "NSS search base for passwd map",
				Optional:    true,
				Computed:    true,
			},
			"nss_base_group": schema.StringAttribute{
				Description: "NSS search base for group map",
				Optional:    true,
				Computed:    true,
			},
			"nss_base_shadow": schema.StringAttribute{
				Description: "NSS search base for shadow map",
				Optional:    true,
				Computed:    true,
			},
			"nss_base_netgroup": schema.StringAttribute{
				Description: "NSS search base for netgroup map",
				Optional:    true,
				Computed:    true,
			},
			"nss_base_sudoers": schema.StringAttribute{
				Description: "NSS search base for sudoers map",
				Optional:    true,
				Computed:    true,
			},
			"nss_initgroups_ignore_users": schema.StringAttribute{
				Description: "Users for which initgroups lookups are skipped",
				Optional:    true,
				Computed:    true,
			},
			"nss_skip_members": schema.BoolAttribute{
				Description: "Return group entries without member attributes",
				Optional:    true,
				Computed:    true,
			},
			"pam_filter": schema.StringAttribute{
				Description: "PAM search filter for retrieving user information",
				Optional:    true,
				Computed:    true,
			},
			"pam_login_attribute": schema.StringAttribute{
				Description: "Attribute used for the user's login name",
				Optional:    true,
				Computed:    true,
			},
			"pam_group_dn": schema.StringAttribute{
				Description: "Required PAM group DN",
				Optional:    true,
				Computed:    true,
			},
			"pam_member_attribute": schema.StringAttribute{
				Description: "Attribute used to test PAM group membership",
				Optional:    true,
				Computed:    true,
			},
			"sudoers_base": schema.StringAttribute{
				Description: "Base DN for sudo LDAP queries",
				Optional:    true,
				Computed:    true,
			},
			"sudoers_search_filter": schema.StringAttribute{
				Description: "LDAP filter for sudo LDAP queries",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"ldap_servers": schema.ListNestedBlock{
				Description: "LDAP server entries",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: "Enable this LDAP server entry",
							Optional:    true,
							Computed:    true,
						},
						"server": schema.StringAttribute{
							Description: "LDAP server hostname or address",
							Optional:    true,
							Computed:    true,
						},
						"port": schema.Int64Attribute{
							Description: "Server port",
							Optional:    true,
							Computed:    true,
						},
						"use_type": schema.StringAttribute{
							Description: "LDAP clients using this server",
							Optional:    true,
							Computed:    true,
						},
						"priority": schema.Int64Attribute{
							Description: "Server priority",
							Optional:    true,
							Computed:    true,
						},
						"ssl_tls_mode": schema.StringAttribute{
							Description: "Per-server TLS mode",
							Optional:    true,
							Computed:    true,
						},
						"retransmit_attempts": schema.Int64Attribute{
							Description: "Per-server retransmit attempts",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"attribute_maps": schema.ListNestedBlock{
				Description: "LDAP attribute mapping overrides",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Description: "Enable this mapping entry",
							Optional:    true,
							Computed:    true,
						},
						"map_name": schema.StringAttribute{
							Description: "Category of mapping override",
							Optional:    true,
							Computed:    true,
						},
						"from": schema.StringAttribute{
							Description: "Original attribute or class name",
							Optional:    true,
							Computed:    true,
						},
						"to": schema.StringAttribute{
							Description: "Replacement attribute, class name, or value",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityLdapProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityLdapProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityLdapProfileResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ldapProfileProps := &openapi.LdapprofilesPutRequestLdapProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "BaseDn", APIField: &ldapProfileProps.BaseDn, TFValue: plan.BaseDn},
		{FieldName: "BindDn", APIField: &ldapProfileProps.BindDn, TFValue: plan.BindDn},
		{FieldName: "BindPassword", APIField: &ldapProfileProps.BindPassword, TFValue: plan.BindPassword},
		{FieldName: "EncryptedBindPassword", APIField: &ldapProfileProps.EncryptedBindPassword, TFValue: plan.EncryptedBindPassword},
		{FieldName: "LdapVersion", APIField: &ldapProfileProps.LdapVersion, TFValue: plan.LdapVersion},
		{FieldName: "SslTlsMode", APIField: &ldapProfileProps.SslTlsMode, TFValue: plan.SslTlsMode},
		{FieldName: "SearchScope", APIField: &ldapProfileProps.SearchScope, TFValue: plan.SearchScope},
		{FieldName: "NssBasePasswd", APIField: &ldapProfileProps.NssBasePasswd, TFValue: plan.NssBasePasswd},
		{FieldName: "NssBaseGroup", APIField: &ldapProfileProps.NssBaseGroup, TFValue: plan.NssBaseGroup},
		{FieldName: "NssBaseShadow", APIField: &ldapProfileProps.NssBaseShadow, TFValue: plan.NssBaseShadow},
		{FieldName: "NssBaseNetgroup", APIField: &ldapProfileProps.NssBaseNetgroup, TFValue: plan.NssBaseNetgroup},
		{FieldName: "NssBaseSudoers", APIField: &ldapProfileProps.NssBaseSudoers, TFValue: plan.NssBaseSudoers},
		{FieldName: "NssInitgroupsIgnoreUsers", APIField: &ldapProfileProps.NssInitgroupsIgnoreUsers, TFValue: plan.NssInitgroupsIgnoreUsers},
		{FieldName: "PamFilter", APIField: &ldapProfileProps.PamFilter, TFValue: plan.PamFilter},
		{FieldName: "PamLoginAttribute", APIField: &ldapProfileProps.PamLoginAttribute, TFValue: plan.PamLoginAttribute},
		{FieldName: "PamGroupDn", APIField: &ldapProfileProps.PamGroupDn, TFValue: plan.PamGroupDn},
		{FieldName: "PamMemberAttribute", APIField: &ldapProfileProps.PamMemberAttribute, TFValue: plan.PamMemberAttribute},
		{FieldName: "SudoersBase", APIField: &ldapProfileProps.SudoersBase, TFValue: plan.SudoersBase},
		{FieldName: "SudoersSearchFilter", APIField: &ldapProfileProps.SudoersSearchFilter, TFValue: plan.SudoersSearchFilter},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ldapProfileProps.Enable, TFValue: plan.Enable},
		{FieldName: "NssSkipMembers", APIField: &ldapProfileProps.NssSkipMembers, TFValue: plan.NssSkipMembers},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ldapProfileTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "DefaultPort", APIField: &ldapProfileProps.DefaultPort, TFValue: config.DefaultPort, IsConfigured: configuredAttrs.IsConfigured("default_port")},
		{FieldName: "SearchTimeLimit", APIField: &ldapProfileProps.SearchTimeLimit, TFValue: config.SearchTimeLimit, IsConfigured: configuredAttrs.IsConfigured("search_time_limit")},
		{FieldName: "BindTimeLimit", APIField: &ldapProfileProps.BindTimeLimit, TFValue: config.BindTimeLimit, IsConfigured: configuredAttrs.IsConfigured("bind_time_limit")},
		{FieldName: "IdleTimeLimit", APIField: &ldapProfileProps.IdleTimeLimit, TFValue: config.IdleTimeLimit, IsConfigured: configuredAttrs.IsConfigured("idle_time_limit")},
		{FieldName: "RetransmitAttempts", APIField: &ldapProfileProps.RetransmitAttempts, TFValue: config.RetransmitAttempts, IsConfigured: configuredAttrs.IsConfigured("retransmit_attempts")},
	})

	// Handle LDAP servers
	if len(plan.LdapServers) > 0 {
		ldapServersConfigMap := utils.BuildIndexedConfigMap(config.LdapServers)
		ldapServers := make([]openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner, len(plan.LdapServers))
		for i, item := range plan.LdapServers {
			server := openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &server.Enabled, TFValue: item.Enabled},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Server", APIField: &server.Server, TFValue: item.Server},
				{FieldName: "UseType", APIField: &server.UseType, TFValue: item.UseType},
				{FieldName: "SslTlsMode", APIField: &server.SslTlsMode, TFValue: item.SslTlsMode},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(item, ldapServersConfigMap, "ldap_servers", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "Port", APIField: &server.Port, TFValue: configItem.Port, IsConfigured: cfg.IsFieldConfigured("port")},
				{FieldName: "Priority", APIField: &server.Priority, TFValue: configItem.Priority, IsConfigured: cfg.IsFieldConfigured("priority")},
				{FieldName: "RetransmitAttempts", APIField: &server.RetransmitAttempts, TFValue: configItem.RetransmitAttempts, IsConfigured: cfg.IsFieldConfigured("retransmit_attempts")},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: item.Index},
			})
			ldapServers[i] = server
		}
		ldapProfileProps.LdapServers = ldapServers
	}

	// Handle attribute maps
	if len(plan.AttributeMaps) > 0 {
		attributeMaps := make([]openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner, len(plan.AttributeMaps))
		for i, item := range plan.AttributeMaps {
			attributeMap := openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &attributeMap.Enabled, TFValue: item.Enabled},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "MapName", APIField: &attributeMap.MapName, TFValue: item.MapName},
				{FieldName: "From", APIField: &attributeMap.From, TFValue: item.From},
				{FieldName: "To", APIField: &attributeMap.To, TFValue: item.To},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &attributeMap.Index, TFValue: item.Index},
			})
			attributeMaps[i] = attributeMap
		}
		ldapProfileProps.AttributeMaps = attributeMaps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ldap_profile", name, *ldapProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LDAP Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ldap_profiles")

	var minState verityLdapProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if profileData, exists := bulkMgr.GetResourceResponse("ldap_profile", name); exists {
			state := populateLdapProfileState(ctx, minState, utils.MergeMissingPlanScalars(profileData, plan, ldapProfileResourceType, r.provCtx.mode), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, ldapProfileResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityLdapProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityLdapProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	profileName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if profileData, exists := r.bulkOpsMgr.GetResourceResponse("ldap_profile", profileName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached LDAP profile data for %s from recent operation", profileName))
			state = populateLdapProfileState(ctx, state, utils.ApplyPostOperationFallback(ctx, profileData), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ldap_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping LDAP profile %s verification – trusting recent successful API operation", profileName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching LDAP profiles for verification of %s", profileName))

	type LdapProfilesResponse struct {
		LdapProfile map[string]interface{} `json:"ldap_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "ldap_profiles", profileName,
		func() (LdapProfilesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch LDAP profiles")
			respAPI, err := r.client.LDAPProfilesAPI.LdapprofilesGet(ctx).Execute()
			if err != nil {
				return LdapProfilesResponse{}, fmt.Errorf("error reading LDAP profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res LdapProfilesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return LdapProfilesResponse{}, fmt.Errorf("failed to decode LDAP profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d LDAP profiles", len(res.LdapProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read LDAP Profile %s", profileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for LDAP profile with name: %s", profileName))

	profileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.LdapProfile,
		profileName,
		func(data interface{}) (string, bool) {
			if profile, ok := data.(map[string]interface{}); ok {
				if name, ok := profile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("LDAP profile with name '%s' not found in API response", profileName))
		resp.State.RemoveResource(ctx)
		return
	}

	profileMap, ok := profileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid LDAP Profile Data",
			fmt.Sprintf("LDAP Profile data is not in expected format for %s", profileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found LDAP profile '%s' under API key '%s'", profileName, actualAPIName))

	state = populateLdapProfileState(ctx, state, utils.ApplyPostOperationFallback(ctx, profileMap), r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityLdapProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityLdapProfileResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get config for nullable field handling
	var config verityLdapProfileResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ldapProfileProps := openapi.LdapprofilesPutRequestLdapProfileValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ldapProfileTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ldapProfileProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.BaseDn, state.BaseDn, func(v *string) { ldapProfileProps.BaseDn = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.BindDn, state.BindDn, func(v *string) { ldapProfileProps.BindDn = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.BindPassword, state.BindPassword, func(v *string) { ldapProfileProps.BindPassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.EncryptedBindPassword, state.EncryptedBindPassword, func(v *string) { ldapProfileProps.EncryptedBindPassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.LdapVersion, state.LdapVersion, func(v *string) { ldapProfileProps.LdapVersion = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SslTlsMode, state.SslTlsMode, func(v *string) { ldapProfileProps.SslTlsMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SearchScope, state.SearchScope, func(v *string) { ldapProfileProps.SearchScope = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NssBasePasswd, state.NssBasePasswd, func(v *string) { ldapProfileProps.NssBasePasswd = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NssBaseGroup, state.NssBaseGroup, func(v *string) { ldapProfileProps.NssBaseGroup = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NssBaseShadow, state.NssBaseShadow, func(v *string) { ldapProfileProps.NssBaseShadow = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NssBaseNetgroup, state.NssBaseNetgroup, func(v *string) { ldapProfileProps.NssBaseNetgroup = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NssBaseSudoers, state.NssBaseSudoers, func(v *string) { ldapProfileProps.NssBaseSudoers = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NssInitgroupsIgnoreUsers, state.NssInitgroupsIgnoreUsers, func(v *string) { ldapProfileProps.NssInitgroupsIgnoreUsers = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PamFilter, state.PamFilter, func(v *string) { ldapProfileProps.PamFilter = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PamLoginAttribute, state.PamLoginAttribute, func(v *string) { ldapProfileProps.PamLoginAttribute = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PamGroupDn, state.PamGroupDn, func(v *string) { ldapProfileProps.PamGroupDn = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PamMemberAttribute, state.PamMemberAttribute, func(v *string) { ldapProfileProps.PamMemberAttribute = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SudoersBase, state.SudoersBase, func(v *string) { ldapProfileProps.SudoersBase = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SudoersSearchFilter, state.SudoersSearchFilter, func(v *string) { ldapProfileProps.SudoersSearchFilter = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ldapProfileProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.NssSkipMembers, state.NssSkipMembers, func(v *bool) { ldapProfileProps.NssSkipMembers = v }, &hasChanges)

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.DefaultPort, state.DefaultPort, configuredAttrs.IsConfigured("default_port"), func(v *openapi.NullableInt64) { ldapProfileProps.DefaultPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.SearchTimeLimit, state.SearchTimeLimit, configuredAttrs.IsConfigured("search_time_limit"), func(v *openapi.NullableInt64) { ldapProfileProps.SearchTimeLimit = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.BindTimeLimit, state.BindTimeLimit, configuredAttrs.IsConfigured("bind_time_limit"), func(v *openapi.NullableInt64) { ldapProfileProps.BindTimeLimit = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.IdleTimeLimit, state.IdleTimeLimit, configuredAttrs.IsConfigured("idle_time_limit"), func(v *openapi.NullableInt64) { ldapProfileProps.IdleTimeLimit = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.RetransmitAttempts, state.RetransmitAttempts, configuredAttrs.IsConfigured("retransmit_attempts"), func(v *openapi.NullableInt64) { ldapProfileProps.RetransmitAttempts = *v }, &hasChanges)

	// Handle LDAP servers
	ldapServersConfigMap := utils.BuildIndexedConfigMap(config.LdapServers)

	ldapServersHandler := utils.IndexedItemHandler[verityLdapProfileLdapServerModel, openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner]{
		CreateNew: func(planItem verityLdapProfileLdapServerModel) openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner {
			server := openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &server.Enabled, TFValue: planItem.Enabled},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Server", APIField: &server.Server, TFValue: planItem.Server},
				{FieldName: "UseType", APIField: &server.UseType, TFValue: planItem.UseType},
				{FieldName: "SslTlsMode", APIField: &server.SslTlsMode, TFValue: planItem.SslTlsMode},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(planItem, ldapServersConfigMap, "ldap_servers", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "Port", APIField: &server.Port, TFValue: configItem.Port, IsConfigured: cfg.IsFieldConfigured("port")},
				{FieldName: "Priority", APIField: &server.Priority, TFValue: configItem.Priority, IsConfigured: cfg.IsFieldConfigured("priority")},
				{FieldName: "RetransmitAttempts", APIField: &server.RetransmitAttempts, TFValue: configItem.RetransmitAttempts, IsConfigured: cfg.IsFieldConfigured("retransmit_attempts")},
			})

			return server
		},
		UpdateExisting: func(planItem verityLdapProfileLdapServerModel, stateItem verityLdapProfileLdapServerModel) (openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner, bool) {
			server := openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enabled, stateItem.Enabled, func(v *bool) { server.Enabled = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.Server, stateItem.Server, func(v *string) { server.Server = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.UseType, stateItem.UseType, func(v *string) { server.UseType = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.SslTlsMode, stateItem.SslTlsMode, func(v *string) { server.SslTlsMode = v }, &fieldChanged)

			// Handle nullable int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(planItem, ldapServersConfigMap, "ldap_servers", configuredAttrs)
			utils.CompareAndSetNullableInt64Field(configItem.Port, stateItem.Port, cfg.IsFieldConfigured("port"), func(v *openapi.NullableInt64) { server.Port = *v }, &fieldChanged)
			utils.CompareAndSetNullableInt64Field(configItem.Priority, stateItem.Priority, cfg.IsFieldConfigured("priority"), func(v *openapi.NullableInt64) { server.Priority = *v }, &fieldChanged)
			utils.CompareAndSetNullableInt64Field(configItem.RetransmitAttempts, stateItem.RetransmitAttempts, cfg.IsFieldConfigured("retransmit_attempts"), func(v *openapi.NullableInt64) { server.RetransmitAttempts = *v }, &fieldChanged)

			return server, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner {
			server := openapi.LdapprofilesPutRequestLdapProfileValueLdapServersInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &server.Index, TFValue: types.Int64Value(index)},
			})
			return server
		},
	}

	changedLdapServers, ldapServersChanged := utils.ProcessIndexedArrayUpdates(plan.LdapServers, state.LdapServers, ldapServersHandler)
	if ldapServersChanged {
		ldapProfileProps.LdapServers = changedLdapServers
		hasChanges = true
	}

	// Handle attribute maps
	attributeMapsHandler := utils.IndexedItemHandler[verityLdapProfileAttributeMapModel, openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner]{
		CreateNew: func(planItem verityLdapProfileAttributeMapModel) openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner {
			attributeMap := openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &attributeMap.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enabled", APIField: &attributeMap.Enabled, TFValue: planItem.Enabled},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "MapName", APIField: &attributeMap.MapName, TFValue: planItem.MapName},
				{FieldName: "From", APIField: &attributeMap.From, TFValue: planItem.From},
				{FieldName: "To", APIField: &attributeMap.To, TFValue: planItem.To},
			})

			return attributeMap
		},
		UpdateExisting: func(planItem verityLdapProfileAttributeMapModel, stateItem verityLdapProfileAttributeMapModel) (openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner, bool) {
			attributeMap := openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &attributeMap.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enabled, stateItem.Enabled, func(v *bool) { attributeMap.Enabled = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.MapName, stateItem.MapName, func(v *string) { attributeMap.MapName = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.From, stateItem.From, func(v *string) { attributeMap.From = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.To, stateItem.To, func(v *string) { attributeMap.To = v }, &fieldChanged)

			return attributeMap, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner {
			attributeMap := openapi.LdapprofilesPutRequestLdapProfileValueAttributeMapsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &attributeMap.Index, TFValue: types.Int64Value(index)},
			})
			return attributeMap
		},
	}

	changedAttributeMaps, attributeMapsChanged := utils.ProcessIndexedArrayUpdates(plan.AttributeMaps, state.AttributeMaps, attributeMapsHandler)
	if attributeMapsChanged {
		ldapProfileProps.AttributeMaps = changedAttributeMaps
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ldap_profile", name, ldapProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LDAP Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ldap_profiles")

	var minState verityLdapProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if profileData, exists := bulkMgr.GetResourceResponse("ldap_profile", name); exists {
			newState := populateLdapProfileState(ctx, minState, utils.MergeMissingPlanScalars(profileData, plan, ldapProfileResourceType, r.provCtx.mode), r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, ldapProfileResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityLdapProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityLdapProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ldap_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("LDAP Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ldap_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityLdapProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateLdapProfileState(ctx context.Context, state verityLdapProfileResourceModel, data map[string]interface{}, mode string) verityLdapProfileResourceModel {
	const resourceType = ldapProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Int fields
	state.DefaultPort = utils.MapInt64WithMode(data, "default_port", resourceType, mode)
	state.SearchTimeLimit = utils.MapInt64WithMode(data, "search_time_limit", resourceType, mode)
	state.BindTimeLimit = utils.MapInt64WithMode(data, "bind_time_limit", resourceType, mode)
	state.IdleTimeLimit = utils.MapInt64WithMode(data, "idle_time_limit", resourceType, mode)
	state.RetransmitAttempts = utils.MapInt64WithMode(data, "retransmit_attempts", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.NssSkipMembers = utils.MapBoolWithMode(data, "nss_skip_members", resourceType, mode)

	// String fields
	state.BaseDn = utils.MapStringWithMode(data, "base_dn", resourceType, mode)
	state.BindDn = utils.MapStringWithMode(data, "bind_dn", resourceType, mode)
	state.BindPassword = utils.MapStringWithMode(data, "bind_password", resourceType, mode)
	state.EncryptedBindPassword = utils.MapStringWithMode(data, "encrypted_bind_password", resourceType, mode)
	state.LdapVersion = utils.MapStringWithMode(data, "ldap_version", resourceType, mode)
	state.SslTlsMode = utils.MapStringWithMode(data, "ssl_tls_mode", resourceType, mode)
	state.SearchScope = utils.MapStringWithMode(data, "search_scope", resourceType, mode)
	state.NssBasePasswd = utils.MapStringWithMode(data, "nss_base_passwd", resourceType, mode)
	state.NssBaseGroup = utils.MapStringWithMode(data, "nss_base_group", resourceType, mode)
	state.NssBaseShadow = utils.MapStringWithMode(data, "nss_base_shadow", resourceType, mode)
	state.NssBaseNetgroup = utils.MapStringWithMode(data, "nss_base_netgroup", resourceType, mode)
	state.NssBaseSudoers = utils.MapStringWithMode(data, "nss_base_sudoers", resourceType, mode)
	state.NssInitgroupsIgnoreUsers = utils.MapStringWithMode(data, "nss_initgroups_ignore_users", resourceType, mode)
	state.PamFilter = utils.MapStringWithMode(data, "pam_filter", resourceType, mode)
	state.PamLoginAttribute = utils.MapStringWithMode(data, "pam_login_attribute", resourceType, mode)
	state.PamGroupDn = utils.MapStringWithMode(data, "pam_group_dn", resourceType, mode)
	state.PamMemberAttribute = utils.MapStringWithMode(data, "pam_member_attribute", resourceType, mode)
	state.SudoersBase = utils.MapStringWithMode(data, "sudoers_base", resourceType, mode)
	state.SudoersSearchFilter = utils.MapStringWithMode(data, "sudoers_search_filter", resourceType, mode)

	// Handle ldap_servers list block
	if utils.FieldAppliesToMode(resourceType, "ldap_servers", mode) {
		if entries, ok := data["ldap_servers"].([]interface{}); ok && len(entries) > 0 {
			var ldapServers []verityLdapProfileLdapServerModel

			for _, entry := range entries {
				item, ok := entry.(map[string]interface{})
				if !ok {
					continue
				}

				ldapServer := verityLdapProfileLdapServerModel{
					Enabled:            utils.MapBoolWithModeNested(item, "enabled", resourceType, "ldap_servers.enabled", mode),
					Server:             utils.MapStringWithModeNested(item, "server", resourceType, "ldap_servers.server", mode),
					Port:               utils.MapInt64WithModeNested(item, "port", resourceType, "ldap_servers.port", mode),
					UseType:            utils.MapStringWithModeNested(item, "use_type", resourceType, "ldap_servers.use_type", mode),
					Priority:           utils.MapInt64WithModeNested(item, "priority", resourceType, "ldap_servers.priority", mode),
					SslTlsMode:         utils.MapStringWithModeNested(item, "ssl_tls_mode", resourceType, "ldap_servers.ssl_tls_mode", mode),
					RetransmitAttempts: utils.MapInt64WithModeNested(item, "retransmit_attempts", resourceType, "ldap_servers.retransmit_attempts", mode),
					Index:              utils.MapInt64WithModeNested(item, "index", resourceType, "ldap_servers.index", mode),
				}

				ldapServers = append(ldapServers, ldapServer)
			}

			state.LdapServers = ldapServers
		} else {
			state.LdapServers = nil
		}
	} else {
		state.LdapServers = nil
	}

	// Handle attribute_maps list block
	if utils.FieldAppliesToMode(resourceType, "attribute_maps", mode) {
		if entries, ok := data["attribute_maps"].([]interface{}); ok && len(entries) > 0 {
			var attributeMaps []verityLdapProfileAttributeMapModel

			for _, entry := range entries {
				item, ok := entry.(map[string]interface{})
				if !ok {
					continue
				}

				attributeMap := verityLdapProfileAttributeMapModel{
					Enabled: utils.MapBoolWithModeNested(item, "enabled", resourceType, "attribute_maps.enabled", mode),
					MapName: utils.MapStringWithModeNested(item, "map_name", resourceType, "attribute_maps.map_name", mode),
					From:    utils.MapStringWithModeNested(item, "from", resourceType, "attribute_maps.from", mode),
					To:      utils.MapStringWithModeNested(item, "to", resourceType, "attribute_maps.to", mode),
					Index:   utils.MapInt64WithModeNested(item, "index", resourceType, "attribute_maps.index", mode),
				}

				attributeMaps = append(attributeMaps, attributeMap)
			}

			state.AttributeMaps = attributeMaps
		} else {
			state.AttributeMaps = nil
		}
	} else {
		state.AttributeMaps = nil
	}

	return state
}

func (r *verityLdapProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityLdapProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = ldapProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"base_dn", "bind_dn", "bind_password", "encrypted_bind_password", "ldap_version", "ssl_tls_mode",
		"search_scope", "nss_base_passwd", "nss_base_group", "nss_base_shadow", "nss_base_netgroup",
		"nss_base_sudoers", "nss_initgroups_ignore_users", "pam_filter", "pam_login_attribute",
		"pam_group_dn", "pam_member_attribute", "sudoers_base", "sudoers_search_filter",
	)

	nullifier.NullifyBools("enable", "nss_skip_members")

	nullifier.NullifyInt64s("default_port", "search_time_limit", "bind_time_limit", "idle_time_limit", "retransmit_attempts")

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "ldap_servers",
		ItemCount:    len(plan.LdapServers),
		StringFields: []string{"server", "use_type", "ssl_tls_mode"},
		BoolFields:   []string{"enabled"},
		Int64Fields:  []string{"port", "priority", "retransmit_attempts", "index"},
	})
	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "attribute_maps",
		ItemCount:    len(plan.AttributeMaps),
		StringFields: []string{"map_name", "from", "to"},
		BoolFields:   []string{"enabled"},
		Int64Fields:  []string{"index"},
	})

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityLdapProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityLdapProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable Int64 fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ldapProfileTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "default_port", ConfigVal: config.DefaultPort, StateVal: state.DefaultPort},
			{AttrName: "search_time_limit", ConfigVal: config.SearchTimeLimit, StateVal: state.SearchTimeLimit},
			{AttrName: "bind_time_limit", ConfigVal: config.BindTimeLimit, StateVal: state.BindTimeLimit},
			{AttrName: "idle_time_limit", ConfigVal: config.IdleTimeLimit, StateVal: state.IdleTimeLimit},
			{AttrName: "retransmit_attempts", ConfigVal: config.RetransmitAttempts, StateVal: state.RetransmitAttempts},
		},
	})

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	for i, configServer := range config.LdapServers {
		serverIndex := configServer.Index.ValueInt64()
		var stateServer *verityLdapProfileLdapServerModel
		for j := range state.LdapServers {
			if state.LdapServers[j].Index.ValueInt64() == serverIndex {
				stateServer = &state.LdapServers[j]
				break
			}
		}

		if stateServer != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "ldap_servers",
				BlockListPath:   "ldap_servers",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: serverIndex, AttrName: "port", ConfigVal: configServer.Port, StateVal: stateServer.Port},
					{BlockIndex: serverIndex, AttrName: "priority", ConfigVal: configServer.Priority, StateVal: stateServer.Priority},
					{BlockIndex: serverIndex, AttrName: "retransmit_attempts", ConfigVal: configServer.RetransmitAttempts, StateVal: stateServer.RetransmitAttempts},
				},
			})
		}
	}
}
