# \TACACSProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TacacsprofilesDelete**](TACACSProfilesAPI.md#TacacsprofilesDelete) | **Delete** /tacacsprofiles | Delete TACACS Profile
[**TacacsprofilesGet**](TACACSProfilesAPI.md#TacacsprofilesGet) | **Get** /tacacsprofiles | Get all TACACS Profiles
[**TacacsprofilesPatch**](TACACSProfilesAPI.md#TacacsprofilesPatch) | **Patch** /tacacsprofiles | Update TACACS Profile
[**TacacsprofilesPut**](TACACSProfilesAPI.md#TacacsprofilesPut) | **Put** /tacacsprofiles | Create TACACS Profile



## TacacsprofilesDelete

> TacacsprofilesDelete(ctx).TacacsProfileName(tacacsProfileName).ChangesetName(changesetName).Execute()

Delete TACACS Profile



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
	tacacsProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.TACACSProfilesAPI.TacacsprofilesDelete(context.Background()).TacacsProfileName(tacacsProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TACACSProfilesAPI.TacacsprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTacacsprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tacacsProfileName** | **[]string** |  | 
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


## TacacsprofilesGet

> TacacsprofilesGet(ctx).TacacsProfileName(tacacsProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all TACACS Profiles



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
	tacacsProfileName := "tacacsProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.TACACSProfilesAPI.TacacsprofilesGet(context.Background()).TacacsProfileName(tacacsProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TACACSProfilesAPI.TacacsprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTacacsprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tacacsProfileName** | **string** |  | 
 **includeData** | **bool** |  | 
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


## TacacsprofilesPatch

> TacacsprofilesPatch(ctx).ChangesetName(changesetName).TacacsprofilesPutRequest(tacacsprofilesPutRequest).Execute()

Update TACACS Profile



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
	tacacsprofilesPutRequest := *openapiclient.NewTacacsprofilesPutRequest() // TacacsprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.TACACSProfilesAPI.TacacsprofilesPatch(context.Background()).ChangesetName(changesetName).TacacsprofilesPutRequest(tacacsprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TACACSProfilesAPI.TacacsprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTacacsprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **tacacsprofilesPutRequest** | [**TacacsprofilesPutRequest**](TacacsprofilesPutRequest.md) |  | 

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


## TacacsprofilesPut

> TacacsprofilesPut(ctx).ChangesetName(changesetName).TacacsprofilesPutRequest(tacacsprofilesPutRequest).Execute()

Create TACACS Profile



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
	tacacsprofilesPutRequest := *openapiclient.NewTacacsprofilesPutRequest() // TacacsprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.TACACSProfilesAPI.TacacsprofilesPut(context.Background()).ChangesetName(changesetName).TacacsprofilesPutRequest(tacacsprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TACACSProfilesAPI.TacacsprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTacacsprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **tacacsprofilesPutRequest** | [**TacacsprofilesPutRequest**](TacacsprofilesPutRequest.md) |  | 

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

