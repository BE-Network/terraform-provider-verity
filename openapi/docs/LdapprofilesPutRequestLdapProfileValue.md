# LdapprofilesPutRequestLdapProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**BaseDn** | Pointer to **string** | Base Distinguished Name to use for LDAP searches | [optional] [default to ""]
**BindDn** | Pointer to **string** | Distinguished Name with which to bind to the LDAP server. Empty value means anonymous bind. | [optional] [default to ""]
**BindPassword** | Pointer to **string** | Credentials with which to bind to the LDAP server. Only used together with Bind DN. | [optional] [default to ""]
**EncryptedBindPassword** | Pointer to **string** | System-generated encrypted version of Bind Password | [optional] [default to ""]
**LdapVersion** | Pointer to **string** | LDAP protocol version | [optional] [default to "3"]
**SslTlsMode** | Pointer to **string** | Global TLS mode for LDAP connections | [optional] [default to "off"]
**DefaultPort** | Pointer to **NullableInt32** | Default LDAP server port (389 for plain/StartTLS, 636 for LDAPS) | [optional] 
**SearchTimeLimit** | Pointer to **NullableInt32** | Search time limit, in seconds | [optional] 
**BindTimeLimit** | Pointer to **NullableInt32** | Bind/connect time limit, in seconds | [optional] 
**IdleTimeLimit** | Pointer to **NullableInt32** | NSS idle connection time limit, in seconds | [optional] 
**RetransmitAttempts** | Pointer to **NullableInt32** | Number of retransmit attempts (0-10) | [optional] 
**SearchScope** | Pointer to **string** | Default LDAP search scope | [optional] [default to "sub"]
**NssBasePasswd** | Pointer to **string** | NSS search base for passwd map | [optional] [default to ""]
**NssBaseGroup** | Pointer to **string** | NSS search base for group map | [optional] [default to ""]
**NssBaseShadow** | Pointer to **string** | NSS search base for shadow map | [optional] [default to ""]
**NssBaseNetgroup** | Pointer to **string** | NSS search base for netgroup map | [optional] [default to ""]
**NssBaseSudoers** | Pointer to **string** | NSS search base for sudoers map | [optional] [default to ""]
**NssInitgroupsIgnoreUsers** | Pointer to **string** | Comma-separated list of users for which initgroups() lookups are skipped | [optional] [default to ""]
**NssSkipMembers** | Pointer to **bool** | If true, the group entry is returned without member attributes | [optional] [default to false]
**PamFilter** | Pointer to **string** | PAM search filter for retrieving user info | [optional] [default to ""]
**PamLoginAttribute** | Pointer to **string** | Attribute used to construct the assertion for the user&#39;s login name | [optional] [default to ""]
**PamGroupDn** | Pointer to **string** | DN of a group a user must belong to for login authorization to succeed | [optional] [default to ""]
**PamMemberAttribute** | Pointer to **string** | Attribute used to test a user&#39;s membership of the PAM group DN | [optional] [default to ""]
**SudoersBase** | Pointer to **string** | Base DN for sudo LDAP queries | [optional] [default to ""]
**SudoersSearchFilter** | Pointer to **string** | LDAP filter used to restrict records returned for sudo LDAP queries | [optional] [default to ""]
**LdapServers** | Pointer to [**[]LdapprofilesPutRequestLdapProfileValueLdapServersInner**](LdapprofilesPutRequestLdapProfileValueLdapServersInner.md) |  | [optional] 
**AttributeMaps** | Pointer to [**[]LdapprofilesPutRequestLdapProfileValueAttributeMapsInner**](LdapprofilesPutRequestLdapProfileValueAttributeMapsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewLdapprofilesPutRequestLdapProfileValue

`func NewLdapprofilesPutRequestLdapProfileValue() *LdapprofilesPutRequestLdapProfileValue`

NewLdapprofilesPutRequestLdapProfileValue instantiates a new LdapprofilesPutRequestLdapProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLdapprofilesPutRequestLdapProfileValueWithDefaults

`func NewLdapprofilesPutRequestLdapProfileValueWithDefaults() *LdapprofilesPutRequestLdapProfileValue`

NewLdapprofilesPutRequestLdapProfileValueWithDefaults instantiates a new LdapprofilesPutRequestLdapProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *LdapprofilesPutRequestLdapProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *LdapprofilesPutRequestLdapProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *LdapprofilesPutRequestLdapProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *LdapprofilesPutRequestLdapProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *LdapprofilesPutRequestLdapProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *LdapprofilesPutRequestLdapProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetBaseDn

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBaseDn() string`

GetBaseDn returns the BaseDn field if non-nil, zero value otherwise.

### GetBaseDnOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBaseDnOk() (*string, bool)`

GetBaseDnOk returns a tuple with the BaseDn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseDn

`func (o *LdapprofilesPutRequestLdapProfileValue) SetBaseDn(v string)`

SetBaseDn sets BaseDn field to given value.

### HasBaseDn

`func (o *LdapprofilesPutRequestLdapProfileValue) HasBaseDn() bool`

HasBaseDn returns a boolean if a field has been set.

### GetBindDn

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBindDn() string`

GetBindDn returns the BindDn field if non-nil, zero value otherwise.

### GetBindDnOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBindDnOk() (*string, bool)`

GetBindDnOk returns a tuple with the BindDn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBindDn

`func (o *LdapprofilesPutRequestLdapProfileValue) SetBindDn(v string)`

SetBindDn sets BindDn field to given value.

### HasBindDn

`func (o *LdapprofilesPutRequestLdapProfileValue) HasBindDn() bool`

HasBindDn returns a boolean if a field has been set.

### GetBindPassword

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBindPassword() string`

GetBindPassword returns the BindPassword field if non-nil, zero value otherwise.

### GetBindPasswordOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBindPasswordOk() (*string, bool)`

GetBindPasswordOk returns a tuple with the BindPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBindPassword

`func (o *LdapprofilesPutRequestLdapProfileValue) SetBindPassword(v string)`

SetBindPassword sets BindPassword field to given value.

### HasBindPassword

`func (o *LdapprofilesPutRequestLdapProfileValue) HasBindPassword() bool`

HasBindPassword returns a boolean if a field has been set.

### GetEncryptedBindPassword

`func (o *LdapprofilesPutRequestLdapProfileValue) GetEncryptedBindPassword() string`

GetEncryptedBindPassword returns the EncryptedBindPassword field if non-nil, zero value otherwise.

### GetEncryptedBindPasswordOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetEncryptedBindPasswordOk() (*string, bool)`

GetEncryptedBindPasswordOk returns a tuple with the EncryptedBindPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptedBindPassword

`func (o *LdapprofilesPutRequestLdapProfileValue) SetEncryptedBindPassword(v string)`

SetEncryptedBindPassword sets EncryptedBindPassword field to given value.

### HasEncryptedBindPassword

`func (o *LdapprofilesPutRequestLdapProfileValue) HasEncryptedBindPassword() bool`

HasEncryptedBindPassword returns a boolean if a field has been set.

### GetLdapVersion

`func (o *LdapprofilesPutRequestLdapProfileValue) GetLdapVersion() string`

GetLdapVersion returns the LdapVersion field if non-nil, zero value otherwise.

### GetLdapVersionOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetLdapVersionOk() (*string, bool)`

GetLdapVersionOk returns a tuple with the LdapVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLdapVersion

`func (o *LdapprofilesPutRequestLdapProfileValue) SetLdapVersion(v string)`

SetLdapVersion sets LdapVersion field to given value.

### HasLdapVersion

`func (o *LdapprofilesPutRequestLdapProfileValue) HasLdapVersion() bool`

HasLdapVersion returns a boolean if a field has been set.

### GetSslTlsMode

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSslTlsMode() string`

GetSslTlsMode returns the SslTlsMode field if non-nil, zero value otherwise.

### GetSslTlsModeOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSslTlsModeOk() (*string, bool)`

GetSslTlsModeOk returns a tuple with the SslTlsMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSslTlsMode

`func (o *LdapprofilesPutRequestLdapProfileValue) SetSslTlsMode(v string)`

SetSslTlsMode sets SslTlsMode field to given value.

### HasSslTlsMode

`func (o *LdapprofilesPutRequestLdapProfileValue) HasSslTlsMode() bool`

HasSslTlsMode returns a boolean if a field has been set.

### GetDefaultPort

`func (o *LdapprofilesPutRequestLdapProfileValue) GetDefaultPort() int32`

GetDefaultPort returns the DefaultPort field if non-nil, zero value otherwise.

### GetDefaultPortOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetDefaultPortOk() (*int32, bool)`

GetDefaultPortOk returns a tuple with the DefaultPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultPort

`func (o *LdapprofilesPutRequestLdapProfileValue) SetDefaultPort(v int32)`

SetDefaultPort sets DefaultPort field to given value.

### HasDefaultPort

`func (o *LdapprofilesPutRequestLdapProfileValue) HasDefaultPort() bool`

HasDefaultPort returns a boolean if a field has been set.

### SetDefaultPortNil

`func (o *LdapprofilesPutRequestLdapProfileValue) SetDefaultPortNil(b bool)`

 SetDefaultPortNil sets the value for DefaultPort to be an explicit nil

### UnsetDefaultPort
`func (o *LdapprofilesPutRequestLdapProfileValue) UnsetDefaultPort()`

UnsetDefaultPort ensures that no value is present for DefaultPort, not even an explicit nil
### GetSearchTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSearchTimeLimit() int32`

GetSearchTimeLimit returns the SearchTimeLimit field if non-nil, zero value otherwise.

### GetSearchTimeLimitOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSearchTimeLimitOk() (*int32, bool)`

GetSearchTimeLimitOk returns a tuple with the SearchTimeLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) SetSearchTimeLimit(v int32)`

SetSearchTimeLimit sets SearchTimeLimit field to given value.

### HasSearchTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) HasSearchTimeLimit() bool`

HasSearchTimeLimit returns a boolean if a field has been set.

### SetSearchTimeLimitNil

`func (o *LdapprofilesPutRequestLdapProfileValue) SetSearchTimeLimitNil(b bool)`

 SetSearchTimeLimitNil sets the value for SearchTimeLimit to be an explicit nil

### UnsetSearchTimeLimit
`func (o *LdapprofilesPutRequestLdapProfileValue) UnsetSearchTimeLimit()`

UnsetSearchTimeLimit ensures that no value is present for SearchTimeLimit, not even an explicit nil
### GetBindTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBindTimeLimit() int32`

GetBindTimeLimit returns the BindTimeLimit field if non-nil, zero value otherwise.

### GetBindTimeLimitOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetBindTimeLimitOk() (*int32, bool)`

GetBindTimeLimitOk returns a tuple with the BindTimeLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBindTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) SetBindTimeLimit(v int32)`

SetBindTimeLimit sets BindTimeLimit field to given value.

### HasBindTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) HasBindTimeLimit() bool`

HasBindTimeLimit returns a boolean if a field has been set.

### SetBindTimeLimitNil

`func (o *LdapprofilesPutRequestLdapProfileValue) SetBindTimeLimitNil(b bool)`

 SetBindTimeLimitNil sets the value for BindTimeLimit to be an explicit nil

### UnsetBindTimeLimit
`func (o *LdapprofilesPutRequestLdapProfileValue) UnsetBindTimeLimit()`

UnsetBindTimeLimit ensures that no value is present for BindTimeLimit, not even an explicit nil
### GetIdleTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) GetIdleTimeLimit() int32`

GetIdleTimeLimit returns the IdleTimeLimit field if non-nil, zero value otherwise.

### GetIdleTimeLimitOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetIdleTimeLimitOk() (*int32, bool)`

GetIdleTimeLimitOk returns a tuple with the IdleTimeLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdleTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) SetIdleTimeLimit(v int32)`

SetIdleTimeLimit sets IdleTimeLimit field to given value.

### HasIdleTimeLimit

`func (o *LdapprofilesPutRequestLdapProfileValue) HasIdleTimeLimit() bool`

HasIdleTimeLimit returns a boolean if a field has been set.

### SetIdleTimeLimitNil

`func (o *LdapprofilesPutRequestLdapProfileValue) SetIdleTimeLimitNil(b bool)`

 SetIdleTimeLimitNil sets the value for IdleTimeLimit to be an explicit nil

### UnsetIdleTimeLimit
`func (o *LdapprofilesPutRequestLdapProfileValue) UnsetIdleTimeLimit()`

UnsetIdleTimeLimit ensures that no value is present for IdleTimeLimit, not even an explicit nil
### GetRetransmitAttempts

`func (o *LdapprofilesPutRequestLdapProfileValue) GetRetransmitAttempts() int32`

GetRetransmitAttempts returns the RetransmitAttempts field if non-nil, zero value otherwise.

### GetRetransmitAttemptsOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetRetransmitAttemptsOk() (*int32, bool)`

GetRetransmitAttemptsOk returns a tuple with the RetransmitAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetransmitAttempts

`func (o *LdapprofilesPutRequestLdapProfileValue) SetRetransmitAttempts(v int32)`

SetRetransmitAttempts sets RetransmitAttempts field to given value.

### HasRetransmitAttempts

`func (o *LdapprofilesPutRequestLdapProfileValue) HasRetransmitAttempts() bool`

HasRetransmitAttempts returns a boolean if a field has been set.

### SetRetransmitAttemptsNil

`func (o *LdapprofilesPutRequestLdapProfileValue) SetRetransmitAttemptsNil(b bool)`

 SetRetransmitAttemptsNil sets the value for RetransmitAttempts to be an explicit nil

### UnsetRetransmitAttempts
`func (o *LdapprofilesPutRequestLdapProfileValue) UnsetRetransmitAttempts()`

UnsetRetransmitAttempts ensures that no value is present for RetransmitAttempts, not even an explicit nil
### GetSearchScope

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSearchScope() string`

GetSearchScope returns the SearchScope field if non-nil, zero value otherwise.

### GetSearchScopeOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSearchScopeOk() (*string, bool)`

GetSearchScopeOk returns a tuple with the SearchScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchScope

`func (o *LdapprofilesPutRequestLdapProfileValue) SetSearchScope(v string)`

SetSearchScope sets SearchScope field to given value.

### HasSearchScope

`func (o *LdapprofilesPutRequestLdapProfileValue) HasSearchScope() bool`

HasSearchScope returns a boolean if a field has been set.

### GetNssBasePasswd

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBasePasswd() string`

GetNssBasePasswd returns the NssBasePasswd field if non-nil, zero value otherwise.

### GetNssBasePasswdOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBasePasswdOk() (*string, bool)`

GetNssBasePasswdOk returns a tuple with the NssBasePasswd field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssBasePasswd

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssBasePasswd(v string)`

SetNssBasePasswd sets NssBasePasswd field to given value.

### HasNssBasePasswd

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssBasePasswd() bool`

HasNssBasePasswd returns a boolean if a field has been set.

### GetNssBaseGroup

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseGroup() string`

GetNssBaseGroup returns the NssBaseGroup field if non-nil, zero value otherwise.

### GetNssBaseGroupOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseGroupOk() (*string, bool)`

GetNssBaseGroupOk returns a tuple with the NssBaseGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssBaseGroup

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssBaseGroup(v string)`

SetNssBaseGroup sets NssBaseGroup field to given value.

### HasNssBaseGroup

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssBaseGroup() bool`

HasNssBaseGroup returns a boolean if a field has been set.

### GetNssBaseShadow

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseShadow() string`

GetNssBaseShadow returns the NssBaseShadow field if non-nil, zero value otherwise.

### GetNssBaseShadowOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseShadowOk() (*string, bool)`

GetNssBaseShadowOk returns a tuple with the NssBaseShadow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssBaseShadow

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssBaseShadow(v string)`

SetNssBaseShadow sets NssBaseShadow field to given value.

### HasNssBaseShadow

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssBaseShadow() bool`

HasNssBaseShadow returns a boolean if a field has been set.

### GetNssBaseNetgroup

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseNetgroup() string`

GetNssBaseNetgroup returns the NssBaseNetgroup field if non-nil, zero value otherwise.

### GetNssBaseNetgroupOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseNetgroupOk() (*string, bool)`

GetNssBaseNetgroupOk returns a tuple with the NssBaseNetgroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssBaseNetgroup

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssBaseNetgroup(v string)`

SetNssBaseNetgroup sets NssBaseNetgroup field to given value.

### HasNssBaseNetgroup

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssBaseNetgroup() bool`

HasNssBaseNetgroup returns a boolean if a field has been set.

### GetNssBaseSudoers

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseSudoers() string`

GetNssBaseSudoers returns the NssBaseSudoers field if non-nil, zero value otherwise.

### GetNssBaseSudoersOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssBaseSudoersOk() (*string, bool)`

GetNssBaseSudoersOk returns a tuple with the NssBaseSudoers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssBaseSudoers

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssBaseSudoers(v string)`

SetNssBaseSudoers sets NssBaseSudoers field to given value.

### HasNssBaseSudoers

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssBaseSudoers() bool`

HasNssBaseSudoers returns a boolean if a field has been set.

### GetNssInitgroupsIgnoreUsers

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssInitgroupsIgnoreUsers() string`

GetNssInitgroupsIgnoreUsers returns the NssInitgroupsIgnoreUsers field if non-nil, zero value otherwise.

### GetNssInitgroupsIgnoreUsersOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssInitgroupsIgnoreUsersOk() (*string, bool)`

GetNssInitgroupsIgnoreUsersOk returns a tuple with the NssInitgroupsIgnoreUsers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssInitgroupsIgnoreUsers

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssInitgroupsIgnoreUsers(v string)`

SetNssInitgroupsIgnoreUsers sets NssInitgroupsIgnoreUsers field to given value.

### HasNssInitgroupsIgnoreUsers

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssInitgroupsIgnoreUsers() bool`

HasNssInitgroupsIgnoreUsers returns a boolean if a field has been set.

### GetNssSkipMembers

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssSkipMembers() bool`

GetNssSkipMembers returns the NssSkipMembers field if non-nil, zero value otherwise.

### GetNssSkipMembersOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetNssSkipMembersOk() (*bool, bool)`

GetNssSkipMembersOk returns a tuple with the NssSkipMembers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNssSkipMembers

`func (o *LdapprofilesPutRequestLdapProfileValue) SetNssSkipMembers(v bool)`

SetNssSkipMembers sets NssSkipMembers field to given value.

### HasNssSkipMembers

`func (o *LdapprofilesPutRequestLdapProfileValue) HasNssSkipMembers() bool`

HasNssSkipMembers returns a boolean if a field has been set.

### GetPamFilter

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamFilter() string`

GetPamFilter returns the PamFilter field if non-nil, zero value otherwise.

### GetPamFilterOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamFilterOk() (*string, bool)`

GetPamFilterOk returns a tuple with the PamFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPamFilter

`func (o *LdapprofilesPutRequestLdapProfileValue) SetPamFilter(v string)`

SetPamFilter sets PamFilter field to given value.

### HasPamFilter

`func (o *LdapprofilesPutRequestLdapProfileValue) HasPamFilter() bool`

HasPamFilter returns a boolean if a field has been set.

### GetPamLoginAttribute

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamLoginAttribute() string`

GetPamLoginAttribute returns the PamLoginAttribute field if non-nil, zero value otherwise.

### GetPamLoginAttributeOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamLoginAttributeOk() (*string, bool)`

GetPamLoginAttributeOk returns a tuple with the PamLoginAttribute field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPamLoginAttribute

`func (o *LdapprofilesPutRequestLdapProfileValue) SetPamLoginAttribute(v string)`

SetPamLoginAttribute sets PamLoginAttribute field to given value.

### HasPamLoginAttribute

`func (o *LdapprofilesPutRequestLdapProfileValue) HasPamLoginAttribute() bool`

HasPamLoginAttribute returns a boolean if a field has been set.

### GetPamGroupDn

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamGroupDn() string`

GetPamGroupDn returns the PamGroupDn field if non-nil, zero value otherwise.

### GetPamGroupDnOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamGroupDnOk() (*string, bool)`

GetPamGroupDnOk returns a tuple with the PamGroupDn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPamGroupDn

`func (o *LdapprofilesPutRequestLdapProfileValue) SetPamGroupDn(v string)`

SetPamGroupDn sets PamGroupDn field to given value.

### HasPamGroupDn

`func (o *LdapprofilesPutRequestLdapProfileValue) HasPamGroupDn() bool`

HasPamGroupDn returns a boolean if a field has been set.

### GetPamMemberAttribute

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamMemberAttribute() string`

GetPamMemberAttribute returns the PamMemberAttribute field if non-nil, zero value otherwise.

### GetPamMemberAttributeOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetPamMemberAttributeOk() (*string, bool)`

GetPamMemberAttributeOk returns a tuple with the PamMemberAttribute field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPamMemberAttribute

`func (o *LdapprofilesPutRequestLdapProfileValue) SetPamMemberAttribute(v string)`

SetPamMemberAttribute sets PamMemberAttribute field to given value.

### HasPamMemberAttribute

`func (o *LdapprofilesPutRequestLdapProfileValue) HasPamMemberAttribute() bool`

HasPamMemberAttribute returns a boolean if a field has been set.

### GetSudoersBase

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSudoersBase() string`

GetSudoersBase returns the SudoersBase field if non-nil, zero value otherwise.

### GetSudoersBaseOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSudoersBaseOk() (*string, bool)`

GetSudoersBaseOk returns a tuple with the SudoersBase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSudoersBase

`func (o *LdapprofilesPutRequestLdapProfileValue) SetSudoersBase(v string)`

SetSudoersBase sets SudoersBase field to given value.

### HasSudoersBase

`func (o *LdapprofilesPutRequestLdapProfileValue) HasSudoersBase() bool`

HasSudoersBase returns a boolean if a field has been set.

### GetSudoersSearchFilter

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSudoersSearchFilter() string`

GetSudoersSearchFilter returns the SudoersSearchFilter field if non-nil, zero value otherwise.

### GetSudoersSearchFilterOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetSudoersSearchFilterOk() (*string, bool)`

GetSudoersSearchFilterOk returns a tuple with the SudoersSearchFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSudoersSearchFilter

`func (o *LdapprofilesPutRequestLdapProfileValue) SetSudoersSearchFilter(v string)`

SetSudoersSearchFilter sets SudoersSearchFilter field to given value.

### HasSudoersSearchFilter

`func (o *LdapprofilesPutRequestLdapProfileValue) HasSudoersSearchFilter() bool`

HasSudoersSearchFilter returns a boolean if a field has been set.

### GetLdapServers

`func (o *LdapprofilesPutRequestLdapProfileValue) GetLdapServers() []LdapprofilesPutRequestLdapProfileValueLdapServersInner`

GetLdapServers returns the LdapServers field if non-nil, zero value otherwise.

### GetLdapServersOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetLdapServersOk() (*[]LdapprofilesPutRequestLdapProfileValueLdapServersInner, bool)`

GetLdapServersOk returns a tuple with the LdapServers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLdapServers

`func (o *LdapprofilesPutRequestLdapProfileValue) SetLdapServers(v []LdapprofilesPutRequestLdapProfileValueLdapServersInner)`

SetLdapServers sets LdapServers field to given value.

### HasLdapServers

`func (o *LdapprofilesPutRequestLdapProfileValue) HasLdapServers() bool`

HasLdapServers returns a boolean if a field has been set.

### GetAttributeMaps

`func (o *LdapprofilesPutRequestLdapProfileValue) GetAttributeMaps() []LdapprofilesPutRequestLdapProfileValueAttributeMapsInner`

GetAttributeMaps returns the AttributeMaps field if non-nil, zero value otherwise.

### GetAttributeMapsOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetAttributeMapsOk() (*[]LdapprofilesPutRequestLdapProfileValueAttributeMapsInner, bool)`

GetAttributeMapsOk returns a tuple with the AttributeMaps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributeMaps

`func (o *LdapprofilesPutRequestLdapProfileValue) SetAttributeMaps(v []LdapprofilesPutRequestLdapProfileValueAttributeMapsInner)`

SetAttributeMaps sets AttributeMaps field to given value.

### HasAttributeMaps

`func (o *LdapprofilesPutRequestLdapProfileValue) HasAttributeMaps() bool`

HasAttributeMaps returns a boolean if a field has been set.

### GetObjectProperties

`func (o *LdapprofilesPutRequestLdapProfileValue) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *LdapprofilesPutRequestLdapProfileValue) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *LdapprofilesPutRequestLdapProfileValue) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *LdapprofilesPutRequestLdapProfileValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


