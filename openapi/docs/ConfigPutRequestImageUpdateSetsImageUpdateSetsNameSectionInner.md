# ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EndpointSetNumName** | Pointer to **string** | The name of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumTargetUpgradeVersion** | Pointer to **string** | The target SW version for member devices of the Endpoint Set | [optional] [default to "unmanaged"]
**EndpointSetNumUniqueIdentifier** | Pointer to **string** | Unique Identifier - not editable | [optional] [default to "17512797893881"]
**EndpointSetNumOnSummary** | Pointer to **bool** | Include on the Summary | [optional] [default to true]
**EndpointSetNumTargetUpgradeVersionTime** | Pointer to **string** | The time to update to the target SW version | [optional] [default to ""]
**EndpointSetNumSubrule1Inverted** | Pointer to **bool** | Subrule 1 Inverted of the Endpoint Set | [optional] [default to false]
**EndpointSetNumSubrule1Type** | Pointer to **string** | Subrule 1 Type of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule1Value** | Pointer to **string** | Subrule 1 Value of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule1ReferencePath** | Pointer to **string** | Subrule 1 Reference Path of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule1ReferencePathRefType** | Pointer to **string** | Object type for endpoint_set_num_subrule_1_reference_path field | [optional] 
**EndpointSetNumSubrule2Inverted** | Pointer to **bool** | Subrule 2 Inverted of the Endpoint Set | [optional] [default to false]
**EndpointSetNumSubrule2Type** | Pointer to **string** | Subrule 2 Type of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule2Value** | Pointer to **string** | Subrule 2 Value of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule2ReferencePath** | Pointer to **string** | Subrule 2 Reference Path of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule2ReferencePathRefType** | Pointer to **string** | Object type for endpoint_set_num_subrule_2_reference_path field | [optional] 
**EndpointSetNumSubrule3Inverted** | Pointer to **bool** | Subrule 3 Inverted of the Endpoint Set | [optional] [default to false]
**EndpointSetNumSubrule3Type** | Pointer to **string** | Subrule 3 Type of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule3Value** | Pointer to **string** | Subrule 3 Value of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule3ReferencePath** | Pointer to **string** | Subrule 3 Reference Path of the Endpoint Set | [optional] [default to ""]
**EndpointSetNumSubrule3ReferencePathRefType** | Pointer to **string** | Object type for endpoint_set_num_subrule_3_reference_path field | [optional] 

## Methods

### NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner

`func NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner() *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner`

NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner instantiates a new ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInnerWithDefaults

`func NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInnerWithDefaults() *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner`

NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInnerWithDefaults instantiates a new ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEndpointSetNumName

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumName() string`

GetEndpointSetNumName returns the EndpointSetNumName field if non-nil, zero value otherwise.

### GetEndpointSetNumNameOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumNameOk() (*string, bool)`

GetEndpointSetNumNameOk returns a tuple with the EndpointSetNumName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumName

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumName(v string)`

SetEndpointSetNumName sets EndpointSetNumName field to given value.

### HasEndpointSetNumName

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumName() bool`

HasEndpointSetNumName returns a boolean if a field has been set.

### GetEndpointSetNumTargetUpgradeVersion

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumTargetUpgradeVersion() string`

GetEndpointSetNumTargetUpgradeVersion returns the EndpointSetNumTargetUpgradeVersion field if non-nil, zero value otherwise.

### GetEndpointSetNumTargetUpgradeVersionOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumTargetUpgradeVersionOk() (*string, bool)`

GetEndpointSetNumTargetUpgradeVersionOk returns a tuple with the EndpointSetNumTargetUpgradeVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumTargetUpgradeVersion

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumTargetUpgradeVersion(v string)`

SetEndpointSetNumTargetUpgradeVersion sets EndpointSetNumTargetUpgradeVersion field to given value.

### HasEndpointSetNumTargetUpgradeVersion

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumTargetUpgradeVersion() bool`

HasEndpointSetNumTargetUpgradeVersion returns a boolean if a field has been set.

### GetEndpointSetNumUniqueIdentifier

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumUniqueIdentifier() string`

GetEndpointSetNumUniqueIdentifier returns the EndpointSetNumUniqueIdentifier field if non-nil, zero value otherwise.

### GetEndpointSetNumUniqueIdentifierOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumUniqueIdentifierOk() (*string, bool)`

GetEndpointSetNumUniqueIdentifierOk returns a tuple with the EndpointSetNumUniqueIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumUniqueIdentifier

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumUniqueIdentifier(v string)`

SetEndpointSetNumUniqueIdentifier sets EndpointSetNumUniqueIdentifier field to given value.

### HasEndpointSetNumUniqueIdentifier

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumUniqueIdentifier() bool`

HasEndpointSetNumUniqueIdentifier returns a boolean if a field has been set.

### GetEndpointSetNumOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumOnSummary() bool`

GetEndpointSetNumOnSummary returns the EndpointSetNumOnSummary field if non-nil, zero value otherwise.

### GetEndpointSetNumOnSummaryOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumOnSummaryOk() (*bool, bool)`

GetEndpointSetNumOnSummaryOk returns a tuple with the EndpointSetNumOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumOnSummary(v bool)`

SetEndpointSetNumOnSummary sets EndpointSetNumOnSummary field to given value.

### HasEndpointSetNumOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumOnSummary() bool`

HasEndpointSetNumOnSummary returns a boolean if a field has been set.

### GetEndpointSetNumTargetUpgradeVersionTime

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumTargetUpgradeVersionTime() string`

GetEndpointSetNumTargetUpgradeVersionTime returns the EndpointSetNumTargetUpgradeVersionTime field if non-nil, zero value otherwise.

### GetEndpointSetNumTargetUpgradeVersionTimeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumTargetUpgradeVersionTimeOk() (*string, bool)`

GetEndpointSetNumTargetUpgradeVersionTimeOk returns a tuple with the EndpointSetNumTargetUpgradeVersionTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumTargetUpgradeVersionTime

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumTargetUpgradeVersionTime(v string)`

SetEndpointSetNumTargetUpgradeVersionTime sets EndpointSetNumTargetUpgradeVersionTime field to given value.

### HasEndpointSetNumTargetUpgradeVersionTime

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumTargetUpgradeVersionTime() bool`

HasEndpointSetNumTargetUpgradeVersionTime returns a boolean if a field has been set.

### GetEndpointSetNumSubrule1Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1Inverted() bool`

GetEndpointSetNumSubrule1Inverted returns the EndpointSetNumSubrule1Inverted field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule1InvertedOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1InvertedOk() (*bool, bool)`

GetEndpointSetNumSubrule1InvertedOk returns a tuple with the EndpointSetNumSubrule1Inverted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule1Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule1Inverted(v bool)`

SetEndpointSetNumSubrule1Inverted sets EndpointSetNumSubrule1Inverted field to given value.

### HasEndpointSetNumSubrule1Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule1Inverted() bool`

HasEndpointSetNumSubrule1Inverted returns a boolean if a field has been set.

### GetEndpointSetNumSubrule1Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1Type() string`

GetEndpointSetNumSubrule1Type returns the EndpointSetNumSubrule1Type field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule1TypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1TypeOk() (*string, bool)`

GetEndpointSetNumSubrule1TypeOk returns a tuple with the EndpointSetNumSubrule1Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule1Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule1Type(v string)`

SetEndpointSetNumSubrule1Type sets EndpointSetNumSubrule1Type field to given value.

### HasEndpointSetNumSubrule1Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule1Type() bool`

HasEndpointSetNumSubrule1Type returns a boolean if a field has been set.

### GetEndpointSetNumSubrule1Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1Value() string`

GetEndpointSetNumSubrule1Value returns the EndpointSetNumSubrule1Value field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule1ValueOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1ValueOk() (*string, bool)`

GetEndpointSetNumSubrule1ValueOk returns a tuple with the EndpointSetNumSubrule1Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule1Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule1Value(v string)`

SetEndpointSetNumSubrule1Value sets EndpointSetNumSubrule1Value field to given value.

### HasEndpointSetNumSubrule1Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule1Value() bool`

HasEndpointSetNumSubrule1Value returns a boolean if a field has been set.

### GetEndpointSetNumSubrule1ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1ReferencePath() string`

GetEndpointSetNumSubrule1ReferencePath returns the EndpointSetNumSubrule1ReferencePath field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule1ReferencePathOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1ReferencePathOk() (*string, bool)`

GetEndpointSetNumSubrule1ReferencePathOk returns a tuple with the EndpointSetNumSubrule1ReferencePath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule1ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule1ReferencePath(v string)`

SetEndpointSetNumSubrule1ReferencePath sets EndpointSetNumSubrule1ReferencePath field to given value.

### HasEndpointSetNumSubrule1ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule1ReferencePath() bool`

HasEndpointSetNumSubrule1ReferencePath returns a boolean if a field has been set.

### GetEndpointSetNumSubrule1ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1ReferencePathRefType() string`

GetEndpointSetNumSubrule1ReferencePathRefType returns the EndpointSetNumSubrule1ReferencePathRefType field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule1ReferencePathRefTypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule1ReferencePathRefTypeOk() (*string, bool)`

GetEndpointSetNumSubrule1ReferencePathRefTypeOk returns a tuple with the EndpointSetNumSubrule1ReferencePathRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule1ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule1ReferencePathRefType(v string)`

SetEndpointSetNumSubrule1ReferencePathRefType sets EndpointSetNumSubrule1ReferencePathRefType field to given value.

### HasEndpointSetNumSubrule1ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule1ReferencePathRefType() bool`

HasEndpointSetNumSubrule1ReferencePathRefType returns a boolean if a field has been set.

### GetEndpointSetNumSubrule2Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2Inverted() bool`

GetEndpointSetNumSubrule2Inverted returns the EndpointSetNumSubrule2Inverted field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule2InvertedOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2InvertedOk() (*bool, bool)`

GetEndpointSetNumSubrule2InvertedOk returns a tuple with the EndpointSetNumSubrule2Inverted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule2Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule2Inverted(v bool)`

SetEndpointSetNumSubrule2Inverted sets EndpointSetNumSubrule2Inverted field to given value.

### HasEndpointSetNumSubrule2Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule2Inverted() bool`

HasEndpointSetNumSubrule2Inverted returns a boolean if a field has been set.

### GetEndpointSetNumSubrule2Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2Type() string`

GetEndpointSetNumSubrule2Type returns the EndpointSetNumSubrule2Type field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule2TypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2TypeOk() (*string, bool)`

GetEndpointSetNumSubrule2TypeOk returns a tuple with the EndpointSetNumSubrule2Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule2Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule2Type(v string)`

SetEndpointSetNumSubrule2Type sets EndpointSetNumSubrule2Type field to given value.

### HasEndpointSetNumSubrule2Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule2Type() bool`

HasEndpointSetNumSubrule2Type returns a boolean if a field has been set.

### GetEndpointSetNumSubrule2Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2Value() string`

GetEndpointSetNumSubrule2Value returns the EndpointSetNumSubrule2Value field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule2ValueOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2ValueOk() (*string, bool)`

GetEndpointSetNumSubrule2ValueOk returns a tuple with the EndpointSetNumSubrule2Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule2Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule2Value(v string)`

SetEndpointSetNumSubrule2Value sets EndpointSetNumSubrule2Value field to given value.

### HasEndpointSetNumSubrule2Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule2Value() bool`

HasEndpointSetNumSubrule2Value returns a boolean if a field has been set.

### GetEndpointSetNumSubrule2ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2ReferencePath() string`

GetEndpointSetNumSubrule2ReferencePath returns the EndpointSetNumSubrule2ReferencePath field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule2ReferencePathOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2ReferencePathOk() (*string, bool)`

GetEndpointSetNumSubrule2ReferencePathOk returns a tuple with the EndpointSetNumSubrule2ReferencePath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule2ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule2ReferencePath(v string)`

SetEndpointSetNumSubrule2ReferencePath sets EndpointSetNumSubrule2ReferencePath field to given value.

### HasEndpointSetNumSubrule2ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule2ReferencePath() bool`

HasEndpointSetNumSubrule2ReferencePath returns a boolean if a field has been set.

### GetEndpointSetNumSubrule2ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2ReferencePathRefType() string`

GetEndpointSetNumSubrule2ReferencePathRefType returns the EndpointSetNumSubrule2ReferencePathRefType field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule2ReferencePathRefTypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule2ReferencePathRefTypeOk() (*string, bool)`

GetEndpointSetNumSubrule2ReferencePathRefTypeOk returns a tuple with the EndpointSetNumSubrule2ReferencePathRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule2ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule2ReferencePathRefType(v string)`

SetEndpointSetNumSubrule2ReferencePathRefType sets EndpointSetNumSubrule2ReferencePathRefType field to given value.

### HasEndpointSetNumSubrule2ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule2ReferencePathRefType() bool`

HasEndpointSetNumSubrule2ReferencePathRefType returns a boolean if a field has been set.

### GetEndpointSetNumSubrule3Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3Inverted() bool`

GetEndpointSetNumSubrule3Inverted returns the EndpointSetNumSubrule3Inverted field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule3InvertedOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3InvertedOk() (*bool, bool)`

GetEndpointSetNumSubrule3InvertedOk returns a tuple with the EndpointSetNumSubrule3Inverted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule3Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule3Inverted(v bool)`

SetEndpointSetNumSubrule3Inverted sets EndpointSetNumSubrule3Inverted field to given value.

### HasEndpointSetNumSubrule3Inverted

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule3Inverted() bool`

HasEndpointSetNumSubrule3Inverted returns a boolean if a field has been set.

### GetEndpointSetNumSubrule3Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3Type() string`

GetEndpointSetNumSubrule3Type returns the EndpointSetNumSubrule3Type field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule3TypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3TypeOk() (*string, bool)`

GetEndpointSetNumSubrule3TypeOk returns a tuple with the EndpointSetNumSubrule3Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule3Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule3Type(v string)`

SetEndpointSetNumSubrule3Type sets EndpointSetNumSubrule3Type field to given value.

### HasEndpointSetNumSubrule3Type

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule3Type() bool`

HasEndpointSetNumSubrule3Type returns a boolean if a field has been set.

### GetEndpointSetNumSubrule3Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3Value() string`

GetEndpointSetNumSubrule3Value returns the EndpointSetNumSubrule3Value field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule3ValueOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3ValueOk() (*string, bool)`

GetEndpointSetNumSubrule3ValueOk returns a tuple with the EndpointSetNumSubrule3Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule3Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule3Value(v string)`

SetEndpointSetNumSubrule3Value sets EndpointSetNumSubrule3Value field to given value.

### HasEndpointSetNumSubrule3Value

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule3Value() bool`

HasEndpointSetNumSubrule3Value returns a boolean if a field has been set.

### GetEndpointSetNumSubrule3ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3ReferencePath() string`

GetEndpointSetNumSubrule3ReferencePath returns the EndpointSetNumSubrule3ReferencePath field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule3ReferencePathOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3ReferencePathOk() (*string, bool)`

GetEndpointSetNumSubrule3ReferencePathOk returns a tuple with the EndpointSetNumSubrule3ReferencePath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule3ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule3ReferencePath(v string)`

SetEndpointSetNumSubrule3ReferencePath sets EndpointSetNumSubrule3ReferencePath field to given value.

### HasEndpointSetNumSubrule3ReferencePath

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule3ReferencePath() bool`

HasEndpointSetNumSubrule3ReferencePath returns a boolean if a field has been set.

### GetEndpointSetNumSubrule3ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3ReferencePathRefType() string`

GetEndpointSetNumSubrule3ReferencePathRefType returns the EndpointSetNumSubrule3ReferencePathRefType field if non-nil, zero value otherwise.

### GetEndpointSetNumSubrule3ReferencePathRefTypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) GetEndpointSetNumSubrule3ReferencePathRefTypeOk() (*string, bool)`

GetEndpointSetNumSubrule3ReferencePathRefTypeOk returns a tuple with the EndpointSetNumSubrule3ReferencePathRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointSetNumSubrule3ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) SetEndpointSetNumSubrule3ReferencePathRefType(v string)`

SetEndpointSetNumSubrule3ReferencePathRefType sets EndpointSetNumSubrule3ReferencePathRefType field to given value.

### HasEndpointSetNumSubrule3ReferencePathRefType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner) HasEndpointSetNumSubrule3ReferencePathRefType() bool`

HasEndpointSetNumSubrule3ReferencePathRefType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


