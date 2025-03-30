# AuthPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Auth** | Pointer to [**AuthPostRequestAuth**](AuthPostRequestAuth.md) |  | [optional] 

## Methods

### NewAuthPostRequest

`func NewAuthPostRequest() *AuthPostRequest`

NewAuthPostRequest instantiates a new AuthPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthPostRequestWithDefaults

`func NewAuthPostRequestWithDefaults() *AuthPostRequest`

NewAuthPostRequestWithDefaults instantiates a new AuthPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuth

`func (o *AuthPostRequest) GetAuth() AuthPostRequestAuth`

GetAuth returns the Auth field if non-nil, zero value otherwise.

### GetAuthOk

`func (o *AuthPostRequest) GetAuthOk() (*AuthPostRequestAuth, bool)`

GetAuthOk returns a tuple with the Auth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuth

`func (o *AuthPostRequest) SetAuth(v AuthPostRequestAuth)`

SetAuth sets Auth field to given value.

### HasAuth

`func (o *AuthPostRequest) HasAuth() bool`

HasAuth returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


