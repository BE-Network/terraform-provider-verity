# AclsPutRequestIpFilterValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Protocol** | Pointer to **string** | Value must be ip/tcp/udp/icmp or a number between 0 and 255 to match packets.  Value IP will match all IP protocols. | [optional] [default to ""]
**Bidirectional** | Pointer to **bool** | If bidirectional is selected, packets will be selected that match the source filters in either the source or destination fields of the packet. | [optional] [default to false]
**SourceIp** | Pointer to **string** | This field matches the source IP address of an IPv4 packet | [optional] [default to ""]
**SourcePortOperator** | Pointer to **string** | This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range. | [optional] [default to ""]
**SourcePort1** | Pointer to **NullableInt32** | This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation. | [optional] 
**SourcePort2** | Pointer to **NullableInt32** | This field will only be used in the range TCP/UDP port value match operation to define the top value in the range. | [optional] 
**DestinationIp** | Pointer to **string** | This field matches the destination IP address of an IPv4 packet. | [optional] [default to ""]
**DestinationPortOperator** | Pointer to **string** | This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range. | [optional] [default to ""]
**DestinationPort1** | Pointer to **NullableInt32** | This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation. | [optional] 
**DestinationPort2** | Pointer to **NullableInt32** | This field will only be used in the range TCP/UDP port value match operation to define the top value in the range. | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewAclsPutRequestIpFilterValue

`func NewAclsPutRequestIpFilterValue() *AclsPutRequestIpFilterValue`

NewAclsPutRequestIpFilterValue instantiates a new AclsPutRequestIpFilterValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAclsPutRequestIpFilterValueWithDefaults

`func NewAclsPutRequestIpFilterValueWithDefaults() *AclsPutRequestIpFilterValue`

NewAclsPutRequestIpFilterValueWithDefaults instantiates a new AclsPutRequestIpFilterValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *AclsPutRequestIpFilterValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AclsPutRequestIpFilterValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AclsPutRequestIpFilterValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AclsPutRequestIpFilterValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *AclsPutRequestIpFilterValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *AclsPutRequestIpFilterValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *AclsPutRequestIpFilterValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *AclsPutRequestIpFilterValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetProtocol

`func (o *AclsPutRequestIpFilterValue) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *AclsPutRequestIpFilterValue) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *AclsPutRequestIpFilterValue) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *AclsPutRequestIpFilterValue) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetBidirectional

`func (o *AclsPutRequestIpFilterValue) GetBidirectional() bool`

GetBidirectional returns the Bidirectional field if non-nil, zero value otherwise.

### GetBidirectionalOk

`func (o *AclsPutRequestIpFilterValue) GetBidirectionalOk() (*bool, bool)`

GetBidirectionalOk returns a tuple with the Bidirectional field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBidirectional

`func (o *AclsPutRequestIpFilterValue) SetBidirectional(v bool)`

SetBidirectional sets Bidirectional field to given value.

### HasBidirectional

`func (o *AclsPutRequestIpFilterValue) HasBidirectional() bool`

HasBidirectional returns a boolean if a field has been set.

### GetSourceIp

`func (o *AclsPutRequestIpFilterValue) GetSourceIp() string`

GetSourceIp returns the SourceIp field if non-nil, zero value otherwise.

### GetSourceIpOk

`func (o *AclsPutRequestIpFilterValue) GetSourceIpOk() (*string, bool)`

GetSourceIpOk returns a tuple with the SourceIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceIp

`func (o *AclsPutRequestIpFilterValue) SetSourceIp(v string)`

SetSourceIp sets SourceIp field to given value.

### HasSourceIp

`func (o *AclsPutRequestIpFilterValue) HasSourceIp() bool`

HasSourceIp returns a boolean if a field has been set.

### GetSourcePortOperator

`func (o *AclsPutRequestIpFilterValue) GetSourcePortOperator() string`

GetSourcePortOperator returns the SourcePortOperator field if non-nil, zero value otherwise.

### GetSourcePortOperatorOk

`func (o *AclsPutRequestIpFilterValue) GetSourcePortOperatorOk() (*string, bool)`

GetSourcePortOperatorOk returns a tuple with the SourcePortOperator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourcePortOperator

`func (o *AclsPutRequestIpFilterValue) SetSourcePortOperator(v string)`

SetSourcePortOperator sets SourcePortOperator field to given value.

### HasSourcePortOperator

`func (o *AclsPutRequestIpFilterValue) HasSourcePortOperator() bool`

HasSourcePortOperator returns a boolean if a field has been set.

### GetSourcePort1

`func (o *AclsPutRequestIpFilterValue) GetSourcePort1() int32`

GetSourcePort1 returns the SourcePort1 field if non-nil, zero value otherwise.

### GetSourcePort1Ok

`func (o *AclsPutRequestIpFilterValue) GetSourcePort1Ok() (*int32, bool)`

GetSourcePort1Ok returns a tuple with the SourcePort1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourcePort1

`func (o *AclsPutRequestIpFilterValue) SetSourcePort1(v int32)`

SetSourcePort1 sets SourcePort1 field to given value.

### HasSourcePort1

`func (o *AclsPutRequestIpFilterValue) HasSourcePort1() bool`

HasSourcePort1 returns a boolean if a field has been set.

### SetSourcePort1Nil

`func (o *AclsPutRequestIpFilterValue) SetSourcePort1Nil(b bool)`

 SetSourcePort1Nil sets the value for SourcePort1 to be an explicit nil

### UnsetSourcePort1
`func (o *AclsPutRequestIpFilterValue) UnsetSourcePort1()`

UnsetSourcePort1 ensures that no value is present for SourcePort1, not even an explicit nil
### GetSourcePort2

`func (o *AclsPutRequestIpFilterValue) GetSourcePort2() int32`

GetSourcePort2 returns the SourcePort2 field if non-nil, zero value otherwise.

### GetSourcePort2Ok

`func (o *AclsPutRequestIpFilterValue) GetSourcePort2Ok() (*int32, bool)`

GetSourcePort2Ok returns a tuple with the SourcePort2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourcePort2

`func (o *AclsPutRequestIpFilterValue) SetSourcePort2(v int32)`

SetSourcePort2 sets SourcePort2 field to given value.

### HasSourcePort2

`func (o *AclsPutRequestIpFilterValue) HasSourcePort2() bool`

HasSourcePort2 returns a boolean if a field has been set.

### SetSourcePort2Nil

`func (o *AclsPutRequestIpFilterValue) SetSourcePort2Nil(b bool)`

 SetSourcePort2Nil sets the value for SourcePort2 to be an explicit nil

### UnsetSourcePort2
`func (o *AclsPutRequestIpFilterValue) UnsetSourcePort2()`

UnsetSourcePort2 ensures that no value is present for SourcePort2, not even an explicit nil
### GetDestinationIp

`func (o *AclsPutRequestIpFilterValue) GetDestinationIp() string`

GetDestinationIp returns the DestinationIp field if non-nil, zero value otherwise.

### GetDestinationIpOk

`func (o *AclsPutRequestIpFilterValue) GetDestinationIpOk() (*string, bool)`

GetDestinationIpOk returns a tuple with the DestinationIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationIp

`func (o *AclsPutRequestIpFilterValue) SetDestinationIp(v string)`

SetDestinationIp sets DestinationIp field to given value.

### HasDestinationIp

`func (o *AclsPutRequestIpFilterValue) HasDestinationIp() bool`

HasDestinationIp returns a boolean if a field has been set.

### GetDestinationPortOperator

`func (o *AclsPutRequestIpFilterValue) GetDestinationPortOperator() string`

GetDestinationPortOperator returns the DestinationPortOperator field if non-nil, zero value otherwise.

### GetDestinationPortOperatorOk

`func (o *AclsPutRequestIpFilterValue) GetDestinationPortOperatorOk() (*string, bool)`

GetDestinationPortOperatorOk returns a tuple with the DestinationPortOperator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationPortOperator

`func (o *AclsPutRequestIpFilterValue) SetDestinationPortOperator(v string)`

SetDestinationPortOperator sets DestinationPortOperator field to given value.

### HasDestinationPortOperator

`func (o *AclsPutRequestIpFilterValue) HasDestinationPortOperator() bool`

HasDestinationPortOperator returns a boolean if a field has been set.

### GetDestinationPort1

`func (o *AclsPutRequestIpFilterValue) GetDestinationPort1() int32`

GetDestinationPort1 returns the DestinationPort1 field if non-nil, zero value otherwise.

### GetDestinationPort1Ok

`func (o *AclsPutRequestIpFilterValue) GetDestinationPort1Ok() (*int32, bool)`

GetDestinationPort1Ok returns a tuple with the DestinationPort1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationPort1

`func (o *AclsPutRequestIpFilterValue) SetDestinationPort1(v int32)`

SetDestinationPort1 sets DestinationPort1 field to given value.

### HasDestinationPort1

`func (o *AclsPutRequestIpFilterValue) HasDestinationPort1() bool`

HasDestinationPort1 returns a boolean if a field has been set.

### SetDestinationPort1Nil

`func (o *AclsPutRequestIpFilterValue) SetDestinationPort1Nil(b bool)`

 SetDestinationPort1Nil sets the value for DestinationPort1 to be an explicit nil

### UnsetDestinationPort1
`func (o *AclsPutRequestIpFilterValue) UnsetDestinationPort1()`

UnsetDestinationPort1 ensures that no value is present for DestinationPort1, not even an explicit nil
### GetDestinationPort2

`func (o *AclsPutRequestIpFilterValue) GetDestinationPort2() int32`

GetDestinationPort2 returns the DestinationPort2 field if non-nil, zero value otherwise.

### GetDestinationPort2Ok

`func (o *AclsPutRequestIpFilterValue) GetDestinationPort2Ok() (*int32, bool)`

GetDestinationPort2Ok returns a tuple with the DestinationPort2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestinationPort2

`func (o *AclsPutRequestIpFilterValue) SetDestinationPort2(v int32)`

SetDestinationPort2 sets DestinationPort2 field to given value.

### HasDestinationPort2

`func (o *AclsPutRequestIpFilterValue) HasDestinationPort2() bool`

HasDestinationPort2 returns a boolean if a field has been set.

### SetDestinationPort2Nil

`func (o *AclsPutRequestIpFilterValue) SetDestinationPort2Nil(b bool)`

 SetDestinationPort2Nil sets the value for DestinationPort2 to be an explicit nil

### UnsetDestinationPort2
`func (o *AclsPutRequestIpFilterValue) UnsetDestinationPort2()`

UnsetDestinationPort2 ensures that no value is present for DestinationPort2, not even an explicit nil
### GetObjectProperties

`func (o *AclsPutRequestIpFilterValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *AclsPutRequestIpFilterValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *AclsPutRequestIpFilterValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *AclsPutRequestIpFilterValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


