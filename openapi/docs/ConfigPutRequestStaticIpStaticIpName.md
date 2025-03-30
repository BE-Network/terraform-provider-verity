# ConfigPutRequestStaticIpStaticIpName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Ip** | Pointer to **string** | IP | [optional] [default to ""]
**Mac** | Pointer to **string** | MAC | [optional] [default to ""]
**Service** | Pointer to **string** | Service | [optional] [default to ""]
**ServiceRefType** | Pointer to **string** | Object type for service field | [optional] 
**AllowedOn** | Pointer to **string** | Allowed On | [optional] [default to "Static Port"]
**Lag** | Pointer to **string** | LAG | [optional] [default to ""]
**LagRefType** | Pointer to **string** | Object type for lag field | [optional] 
**Switch** | Pointer to **string** |  | [optional] [default to ""]
**SwitchRefType** | Pointer to **string** | Object type for switch field | [optional] 
**Port** | Pointer to **string** | Port | [optional] [default to ""]
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewConfigPutRequestStaticIpStaticIpName

`func NewConfigPutRequestStaticIpStaticIpName() *ConfigPutRequestStaticIpStaticIpName`

NewConfigPutRequestStaticIpStaticIpName instantiates a new ConfigPutRequestStaticIpStaticIpName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestStaticIpStaticIpNameWithDefaults

`func NewConfigPutRequestStaticIpStaticIpNameWithDefaults() *ConfigPutRequestStaticIpStaticIpName`

NewConfigPutRequestStaticIpStaticIpNameWithDefaults instantiates a new ConfigPutRequestStaticIpStaticIpName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestStaticIpStaticIpName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestStaticIpStaticIpName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestStaticIpStaticIpName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestStaticIpStaticIpName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestStaticIpStaticIpName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestStaticIpStaticIpName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIp

`func (o *ConfigPutRequestStaticIpStaticIpName) GetIp() string`

GetIp returns the Ip field if non-nil, zero value otherwise.

### GetIpOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetIpOk() (*string, bool)`

GetIpOk returns a tuple with the Ip field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIp

`func (o *ConfigPutRequestStaticIpStaticIpName) SetIp(v string)`

SetIp sets Ip field to given value.

### HasIp

`func (o *ConfigPutRequestStaticIpStaticIpName) HasIp() bool`

HasIp returns a boolean if a field has been set.

### GetMac

`func (o *ConfigPutRequestStaticIpStaticIpName) GetMac() string`

GetMac returns the Mac field if non-nil, zero value otherwise.

### GetMacOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetMacOk() (*string, bool)`

GetMacOk returns a tuple with the Mac field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMac

`func (o *ConfigPutRequestStaticIpStaticIpName) SetMac(v string)`

SetMac sets Mac field to given value.

### HasMac

`func (o *ConfigPutRequestStaticIpStaticIpName) HasMac() bool`

HasMac returns a boolean if a field has been set.

### GetService

`func (o *ConfigPutRequestStaticIpStaticIpName) GetService() string`

GetService returns the Service field if non-nil, zero value otherwise.

### GetServiceOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetServiceOk() (*string, bool)`

GetServiceOk returns a tuple with the Service field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetService

`func (o *ConfigPutRequestStaticIpStaticIpName) SetService(v string)`

SetService sets Service field to given value.

### HasService

`func (o *ConfigPutRequestStaticIpStaticIpName) HasService() bool`

HasService returns a boolean if a field has been set.

### GetServiceRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) GetServiceRefType() string`

GetServiceRefType returns the ServiceRefType field if non-nil, zero value otherwise.

### GetServiceRefTypeOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetServiceRefTypeOk() (*string, bool)`

GetServiceRefTypeOk returns a tuple with the ServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) SetServiceRefType(v string)`

SetServiceRefType sets ServiceRefType field to given value.

### HasServiceRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) HasServiceRefType() bool`

HasServiceRefType returns a boolean if a field has been set.

### GetAllowedOn

`func (o *ConfigPutRequestStaticIpStaticIpName) GetAllowedOn() string`

GetAllowedOn returns the AllowedOn field if non-nil, zero value otherwise.

### GetAllowedOnOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetAllowedOnOk() (*string, bool)`

GetAllowedOnOk returns a tuple with the AllowedOn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedOn

`func (o *ConfigPutRequestStaticIpStaticIpName) SetAllowedOn(v string)`

SetAllowedOn sets AllowedOn field to given value.

### HasAllowedOn

`func (o *ConfigPutRequestStaticIpStaticIpName) HasAllowedOn() bool`

HasAllowedOn returns a boolean if a field has been set.

### GetLag

`func (o *ConfigPutRequestStaticIpStaticIpName) GetLag() string`

GetLag returns the Lag field if non-nil, zero value otherwise.

### GetLagOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetLagOk() (*string, bool)`

GetLagOk returns a tuple with the Lag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLag

`func (o *ConfigPutRequestStaticIpStaticIpName) SetLag(v string)`

SetLag sets Lag field to given value.

### HasLag

`func (o *ConfigPutRequestStaticIpStaticIpName) HasLag() bool`

HasLag returns a boolean if a field has been set.

### GetLagRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) GetLagRefType() string`

GetLagRefType returns the LagRefType field if non-nil, zero value otherwise.

### GetLagRefTypeOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetLagRefTypeOk() (*string, bool)`

GetLagRefTypeOk returns a tuple with the LagRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLagRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) SetLagRefType(v string)`

SetLagRefType sets LagRefType field to given value.

### HasLagRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) HasLagRefType() bool`

HasLagRefType returns a boolean if a field has been set.

### GetSwitch

`func (o *ConfigPutRequestStaticIpStaticIpName) GetSwitch() string`

GetSwitch returns the Switch field if non-nil, zero value otherwise.

### GetSwitchOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetSwitchOk() (*string, bool)`

GetSwitchOk returns a tuple with the Switch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitch

`func (o *ConfigPutRequestStaticIpStaticIpName) SetSwitch(v string)`

SetSwitch sets Switch field to given value.

### HasSwitch

`func (o *ConfigPutRequestStaticIpStaticIpName) HasSwitch() bool`

HasSwitch returns a boolean if a field has been set.

### GetSwitchRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) GetSwitchRefType() string`

GetSwitchRefType returns the SwitchRefType field if non-nil, zero value otherwise.

### GetSwitchRefTypeOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetSwitchRefTypeOk() (*string, bool)`

GetSwitchRefTypeOk returns a tuple with the SwitchRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) SetSwitchRefType(v string)`

SetSwitchRefType sets SwitchRefType field to given value.

### HasSwitchRefType

`func (o *ConfigPutRequestStaticIpStaticIpName) HasSwitchRefType() bool`

HasSwitchRefType returns a boolean if a field has been set.

### GetPort

`func (o *ConfigPutRequestStaticIpStaticIpName) GetPort() string`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetPortOk() (*string, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ConfigPutRequestStaticIpStaticIpName) SetPort(v string)`

SetPort sets Port field to given value.

### HasPort

`func (o *ConfigPutRequestStaticIpStaticIpName) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestStaticIpStaticIpName) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestStaticIpStaticIpName) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestStaticIpStaticIpName) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestStaticIpStaticIpName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


