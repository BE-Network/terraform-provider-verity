# DeviceaaaprofilesPutRequestDeviceAaaProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**FailThrough** | Pointer to **bool** | When enabled, authentication continues to access each server in the method list if an authentication request fails on one server | [optional] [default to true]
**LoginDefault** | Pointer to [**[]DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner**](DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner.md) |  | [optional] 
**TacacsProfile** | Pointer to **string** | TACACS+ profile for authentication | [optional] [default to ""]
**TacacsProfileRefType** | Pointer to **string** | Object type for tacacs_profile field | [optional] 
**LdapProfile** | Pointer to **string** | LDAP profile for authentication | [optional] [default to ""]
**LdapProfileRefType** | Pointer to **string** | Object type for ldap_profile field | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewDeviceaaaprofilesPutRequestDeviceAaaProfileValue

`func NewDeviceaaaprofilesPutRequestDeviceAaaProfileValue() *DeviceaaaprofilesPutRequestDeviceAaaProfileValue`

NewDeviceaaaprofilesPutRequestDeviceAaaProfileValue instantiates a new DeviceaaaprofilesPutRequestDeviceAaaProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceaaaprofilesPutRequestDeviceAaaProfileValueWithDefaults

`func NewDeviceaaaprofilesPutRequestDeviceAaaProfileValueWithDefaults() *DeviceaaaprofilesPutRequestDeviceAaaProfileValue`

NewDeviceaaaprofilesPutRequestDeviceAaaProfileValueWithDefaults instantiates a new DeviceaaaprofilesPutRequestDeviceAaaProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetFailThrough

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetFailThrough() bool`

GetFailThrough returns the FailThrough field if non-nil, zero value otherwise.

### GetFailThroughOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetFailThroughOk() (*bool, bool)`

GetFailThroughOk returns a tuple with the FailThrough field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailThrough

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetFailThrough(v bool)`

SetFailThrough sets FailThrough field to given value.

### HasFailThrough

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasFailThrough() bool`

HasFailThrough returns a boolean if a field has been set.

### GetLoginDefault

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetLoginDefault() []DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner`

GetLoginDefault returns the LoginDefault field if non-nil, zero value otherwise.

### GetLoginDefaultOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetLoginDefaultOk() (*[]DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner, bool)`

GetLoginDefaultOk returns a tuple with the LoginDefault field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginDefault

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetLoginDefault(v []DeviceaaaprofilesPutRequestDeviceAaaProfileValueLoginDefaultInner)`

SetLoginDefault sets LoginDefault field to given value.

### HasLoginDefault

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasLoginDefault() bool`

HasLoginDefault returns a boolean if a field has been set.

### GetTacacsProfile

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetTacacsProfile() string`

GetTacacsProfile returns the TacacsProfile field if non-nil, zero value otherwise.

### GetTacacsProfileOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetTacacsProfileOk() (*string, bool)`

GetTacacsProfileOk returns a tuple with the TacacsProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTacacsProfile

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetTacacsProfile(v string)`

SetTacacsProfile sets TacacsProfile field to given value.

### HasTacacsProfile

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasTacacsProfile() bool`

HasTacacsProfile returns a boolean if a field has been set.

### GetTacacsProfileRefType

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetTacacsProfileRefType() string`

GetTacacsProfileRefType returns the TacacsProfileRefType field if non-nil, zero value otherwise.

### GetTacacsProfileRefTypeOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetTacacsProfileRefTypeOk() (*string, bool)`

GetTacacsProfileRefTypeOk returns a tuple with the TacacsProfileRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTacacsProfileRefType

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetTacacsProfileRefType(v string)`

SetTacacsProfileRefType sets TacacsProfileRefType field to given value.

### HasTacacsProfileRefType

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasTacacsProfileRefType() bool`

HasTacacsProfileRefType returns a boolean if a field has been set.

### GetLdapProfile

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetLdapProfile() string`

GetLdapProfile returns the LdapProfile field if non-nil, zero value otherwise.

### GetLdapProfileOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetLdapProfileOk() (*string, bool)`

GetLdapProfileOk returns a tuple with the LdapProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLdapProfile

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetLdapProfile(v string)`

SetLdapProfile sets LdapProfile field to given value.

### HasLdapProfile

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasLdapProfile() bool`

HasLdapProfile returns a boolean if a field has been set.

### GetLdapProfileRefType

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetLdapProfileRefType() string`

GetLdapProfileRefType returns the LdapProfileRefType field if non-nil, zero value otherwise.

### GetLdapProfileRefTypeOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetLdapProfileRefTypeOk() (*string, bool)`

GetLdapProfileRefTypeOk returns a tuple with the LdapProfileRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLdapProfileRefType

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetLdapProfileRefType(v string)`

SetLdapProfileRefType sets LdapProfileRefType field to given value.

### HasLdapProfileRefType

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasLdapProfileRefType() bool`

HasLdapProfileRefType returns a boolean if a field has been set.

### GetObjectProperties

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *DeviceaaaprofilesPutRequestDeviceAaaProfileValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


