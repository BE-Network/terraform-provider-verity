# ChangesetsPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Commit** | Pointer to **bool** | Create or commit the changeset | [optional] [default to false]
**ChangesetName** | **string** | Changeset name to create or commit | 

## Methods

### NewChangesetsPostRequest

`func NewChangesetsPostRequest(changesetName string, ) *ChangesetsPostRequest`

NewChangesetsPostRequest instantiates a new ChangesetsPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChangesetsPostRequestWithDefaults

`func NewChangesetsPostRequestWithDefaults() *ChangesetsPostRequest`

NewChangesetsPostRequestWithDefaults instantiates a new ChangesetsPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommit

`func (o *ChangesetsPostRequest) GetCommit() bool`

GetCommit returns the Commit field if non-nil, zero value otherwise.

### GetCommitOk

`func (o *ChangesetsPostRequest) GetCommitOk() (*bool, bool)`

GetCommitOk returns a tuple with the Commit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommit

`func (o *ChangesetsPostRequest) SetCommit(v bool)`

SetCommit sets Commit field to given value.

### HasCommit

`func (o *ChangesetsPostRequest) HasCommit() bool`

HasCommit returns a boolean if a field has been set.

### GetChangesetName

`func (o *ChangesetsPostRequest) GetChangesetName() string`

GetChangesetName returns the ChangesetName field if non-nil, zero value otherwise.

### GetChangesetNameOk

`func (o *ChangesetsPostRequest) GetChangesetNameOk() (*string, bool)`

GetChangesetNameOk returns a tuple with the ChangesetName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChangesetName

`func (o *ChangesetsPostRequest) SetChangesetName(v string)`

SetChangesetName sets ChangesetName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


