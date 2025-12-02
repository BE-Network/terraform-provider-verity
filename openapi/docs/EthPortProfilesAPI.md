# \EthPortProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EthportprofilesDelete**](EthPortProfilesAPI.md#EthportprofilesDelete) | **Delete** /ethportprofiles | Delete Eth-Port Profile
[**EthportprofilesGet**](EthPortProfilesAPI.md#EthportprofilesGet) | **Get** /ethportprofiles | Get all Eth-Port Profiles
[**EthportprofilesPatch**](EthPortProfilesAPI.md#EthportprofilesPatch) | **Patch** /ethportprofiles | Update Eth-Port Profile
[**EthportprofilesPut**](EthPortProfilesAPI.md#EthportprofilesPut) | **Put** /ethportprofiles | Create Eth-Port Profile



## EthportprofilesDelete

> EthportprofilesDelete(ctx).ProfileName(profileName).ChangesetName(changesetName).Execute()

Delete Eth-Port Profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	profileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortProfilesAPI.EthportprofilesDelete(context.Background()).ProfileName(profileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortProfilesAPI.EthportprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **profileName** | **[]string** |  | 
 **changesetName** | **string** |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EthportprofilesGet

> EthportprofilesGet(ctx).ProfileName(profileName).IncludeData(includeData).Execute()

Get all Eth-Port Profiles



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	profileName := "profileName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortProfilesAPI.EthportprofilesGet(context.Background()).ProfileName(profileName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortProfilesAPI.EthportprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **profileName** | **string** |  | 
 **includeData** | **bool** |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EthportprofilesPatch

> EthportprofilesPatch(ctx).ChangesetName(changesetName).EthportprofilesPutRequest(ethportprofilesPutRequest).Execute()

Update Eth-Port Profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	changesetName := "changesetName_example" // string |  (optional)
	ethportprofilesPutRequest := *openapiclient.NewEthportprofilesPutRequest() // EthportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortProfilesAPI.EthportprofilesPatch(context.Background()).ChangesetName(changesetName).EthportprofilesPutRequest(ethportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortProfilesAPI.EthportprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ethportprofilesPutRequest** | [**EthportprofilesPutRequest**](EthportprofilesPutRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EthportprofilesPut

> EthportprofilesPut(ctx).ChangesetName(changesetName).EthportprofilesPutRequest(ethportprofilesPutRequest).Execute()

Create Eth-Port Profile



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	changesetName := "changesetName_example" // string |  (optional)
	ethportprofilesPutRequest := *openapiclient.NewEthportprofilesPutRequest() // EthportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortProfilesAPI.EthportprofilesPut(context.Background()).ChangesetName(changesetName).EthportprofilesPutRequest(ethportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortProfilesAPI.EthportprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ethportprofilesPutRequest** | [**EthportprofilesPutRequest**](EthportprofilesPutRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

