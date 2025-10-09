# PolicybasedroutingPutRequestPbRoutingValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Policy** | Pointer to [**[]PolicybasedroutingPutRequestPbRoutingValuePolicyInner**](PolicybasedroutingPutRequestPbRoutingValuePolicyInner.md) |  | [optional] 

## Methods

### NewPolicybasedroutingPutRequestPbRoutingValue

`func NewPolicybasedroutingPutRequestPbRoutingValue() *PolicybasedroutingPutRequestPbRoutingValue`

NewPolicybasedroutingPutRequestPbRoutingValue instantiates a new PolicybasedroutingPutRequestPbRoutingValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPolicybasedroutingPutRequestPbRoutingValueWithDefaults

`func NewPolicybasedroutingPutRequestPbRoutingValueWithDefaults() *PolicybasedroutingPutRequestPbRoutingValue`

NewPolicybasedroutingPutRequestPbRoutingValueWithDefaults instantiates a new PolicybasedroutingPutRequestPbRoutingValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PolicybasedroutingPutRequestPbRoutingValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PolicybasedroutingPutRequestPbRoutingValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PolicybasedroutingPutRequestPbRoutingValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PolicybasedroutingPutRequestPbRoutingValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PolicybasedroutingPutRequestPbRoutingValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PolicybasedroutingPutRequestPbRoutingValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PolicybasedroutingPutRequestPbRoutingValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PolicybasedroutingPutRequestPbRoutingValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPolicy

`func (o *PolicybasedroutingPutRequestPbRoutingValue) GetPolicy() []PolicybasedroutingPutRequestPbRoutingValuePolicyInner`

GetPolicy returns the Policy field if non-nil, zero value otherwise.

### GetPolicyOk

`func (o *PolicybasedroutingPutRequestPbRoutingValue) GetPolicyOk() (*[]PolicybasedroutingPutRequestPbRoutingValuePolicyInner, bool)`

GetPolicyOk returns a tuple with the Policy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicy

`func (o *PolicybasedroutingPutRequestPbRoutingValue) SetPolicy(v []PolicybasedroutingPutRequestPbRoutingValuePolicyInner)`

SetPolicy sets Policy field to given value.

### HasPolicy

`func (o *PolicybasedroutingPutRequestPbRoutingValue) HasPolicy() bool`

HasPolicy returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


