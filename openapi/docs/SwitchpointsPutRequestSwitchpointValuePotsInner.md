# SwitchpointsPutRequestSwitchpointValuePotsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PotsNumEnable** | Pointer to **bool** | Enable POTS port | [optional] [default to false]
**PotsNumUri** | Pointer to **string** | Specific telephone extension for SIP for POTS port | [optional] [default to ""]
**PotsNumUsername** | Pointer to **string** | SIP username used for authentication for POTS port | [optional] [default to ""]
**PotsNumPassword** | Pointer to **string** | SIP password used for authentication for POTS port | [optional] [default to ""]
**PotsNumCallerId** | Pointer to **string** | ASCII string defining the user for the Caller ID display for POTS port | [optional] [default to ""]
**PotsNumHotLine** | Pointer to **string** | URI of line to autodial upon off-hook for POTS port | [optional] [default to ""]
**PotsNumPasswordEncrypted** | Pointer to **string** | SIP password used for authentication for POTS port | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewSwitchpointsPutRequestSwitchpointValuePotsInner

`func NewSwitchpointsPutRequestSwitchpointValuePotsInner() *SwitchpointsPutRequestSwitchpointValuePotsInner`

NewSwitchpointsPutRequestSwitchpointValuePotsInner instantiates a new SwitchpointsPutRequestSwitchpointValuePotsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsPutRequestSwitchpointValuePotsInnerWithDefaults

`func NewSwitchpointsPutRequestSwitchpointValuePotsInnerWithDefaults() *SwitchpointsPutRequestSwitchpointValuePotsInner`

NewSwitchpointsPutRequestSwitchpointValuePotsInnerWithDefaults instantiates a new SwitchpointsPutRequestSwitchpointValuePotsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPotsNumEnable

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumEnable() bool`

GetPotsNumEnable returns the PotsNumEnable field if non-nil, zero value otherwise.

### GetPotsNumEnableOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumEnableOk() (*bool, bool)`

GetPotsNumEnableOk returns a tuple with the PotsNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumEnable

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumEnable(v bool)`

SetPotsNumEnable sets PotsNumEnable field to given value.

### HasPotsNumEnable

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumEnable() bool`

HasPotsNumEnable returns a boolean if a field has been set.

### GetPotsNumUri

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumUri() string`

GetPotsNumUri returns the PotsNumUri field if non-nil, zero value otherwise.

### GetPotsNumUriOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumUriOk() (*string, bool)`

GetPotsNumUriOk returns a tuple with the PotsNumUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumUri

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumUri(v string)`

SetPotsNumUri sets PotsNumUri field to given value.

### HasPotsNumUri

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumUri() bool`

HasPotsNumUri returns a boolean if a field has been set.

### GetPotsNumUsername

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumUsername() string`

GetPotsNumUsername returns the PotsNumUsername field if non-nil, zero value otherwise.

### GetPotsNumUsernameOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumUsernameOk() (*string, bool)`

GetPotsNumUsernameOk returns a tuple with the PotsNumUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumUsername

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumUsername(v string)`

SetPotsNumUsername sets PotsNumUsername field to given value.

### HasPotsNumUsername

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumUsername() bool`

HasPotsNumUsername returns a boolean if a field has been set.

### GetPotsNumPassword

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumPassword() string`

GetPotsNumPassword returns the PotsNumPassword field if non-nil, zero value otherwise.

### GetPotsNumPasswordOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumPasswordOk() (*string, bool)`

GetPotsNumPasswordOk returns a tuple with the PotsNumPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumPassword

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumPassword(v string)`

SetPotsNumPassword sets PotsNumPassword field to given value.

### HasPotsNumPassword

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumPassword() bool`

HasPotsNumPassword returns a boolean if a field has been set.

### GetPotsNumCallerId

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumCallerId() string`

GetPotsNumCallerId returns the PotsNumCallerId field if non-nil, zero value otherwise.

### GetPotsNumCallerIdOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumCallerIdOk() (*string, bool)`

GetPotsNumCallerIdOk returns a tuple with the PotsNumCallerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumCallerId

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumCallerId(v string)`

SetPotsNumCallerId sets PotsNumCallerId field to given value.

### HasPotsNumCallerId

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumCallerId() bool`

HasPotsNumCallerId returns a boolean if a field has been set.

### GetPotsNumHotLine

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumHotLine() string`

GetPotsNumHotLine returns the PotsNumHotLine field if non-nil, zero value otherwise.

### GetPotsNumHotLineOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumHotLineOk() (*string, bool)`

GetPotsNumHotLineOk returns a tuple with the PotsNumHotLine field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumHotLine

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumHotLine(v string)`

SetPotsNumHotLine sets PotsNumHotLine field to given value.

### HasPotsNumHotLine

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumHotLine() bool`

HasPotsNumHotLine returns a boolean if a field has been set.

### GetPotsNumPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumPasswordEncrypted() string`

GetPotsNumPasswordEncrypted returns the PotsNumPasswordEncrypted field if non-nil, zero value otherwise.

### GetPotsNumPasswordEncryptedOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetPotsNumPasswordEncryptedOk() (*string, bool)`

GetPotsNumPasswordEncryptedOk returns a tuple with the PotsNumPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPotsNumPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetPotsNumPasswordEncrypted(v string)`

SetPotsNumPasswordEncrypted sets PotsNumPasswordEncrypted field to given value.

### HasPotsNumPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasPotsNumPasswordEncrypted() bool`

HasPotsNumPasswordEncrypted returns a boolean if a field has been set.

### GetIndex

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *SwitchpointsPutRequestSwitchpointValuePotsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


