# \ExtendedCommunityListsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ExtendedcommunitylistsDelete**](ExtendedCommunityListsAPI.md#ExtendedcommunitylistsDelete) | **Delete** /extendedcommunitylists | Delete Extended Community List
[**ExtendedcommunitylistsGet**](ExtendedCommunityListsAPI.md#ExtendedcommunitylistsGet) | **Get** /extendedcommunitylists | Get all Extended Community Lists
[**ExtendedcommunitylistsPatch**](ExtendedCommunityListsAPI.md#ExtendedcommunitylistsPatch) | **Patch** /extendedcommunitylists | Update Extended Community List
[**ExtendedcommunitylistsPut**](ExtendedCommunityListsAPI.md#ExtendedcommunitylistsPut) | **Put** /extendedcommunitylists | Create Extended Community List



## ExtendedcommunitylistsDelete

> ExtendedcommunitylistsDelete(ctx).ExtendedCommunityListName(extendedCommunityListName).ChangesetName(changesetName).Execute()

Delete Extended Community List



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
	extendedCommunityListName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ExtendedCommunityListsAPI.ExtendedcommunitylistsDelete(context.Background()).ExtendedCommunityListName(extendedCommunityListName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtendedCommunityListsAPI.ExtendedcommunitylistsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExtendedcommunitylistsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **extendedCommunityListName** | **[]string** |  | 
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


## ExtendedcommunitylistsGet

> ExtendedcommunitylistsGet(ctx).ExtendedCommunityListName(extendedCommunityListName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Extended Community Lists



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
	extendedCommunityListName := "extendedCommunityListName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(context.Background()).ExtendedCommunityListName(extendedCommunityListName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtendedCommunityListsAPI.ExtendedcommunitylistsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExtendedcommunitylistsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **extendedCommunityListName** | **string** |  | 
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


## ExtendedcommunitylistsPatch

> ExtendedcommunitylistsPatch(ctx).ChangesetName(changesetName).ExtendedcommunitylistsPutRequest(extendedcommunitylistsPutRequest).Execute()

Update Extended Community List



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
	extendedcommunitylistsPutRequest := *openapiclient.NewExtendedcommunitylistsPutRequest() // ExtendedcommunitylistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ExtendedCommunityListsAPI.ExtendedcommunitylistsPatch(context.Background()).ChangesetName(changesetName).ExtendedcommunitylistsPutRequest(extendedcommunitylistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtendedCommunityListsAPI.ExtendedcommunitylistsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExtendedcommunitylistsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **extendedcommunitylistsPutRequest** | [**ExtendedcommunitylistsPutRequest**](ExtendedcommunitylistsPutRequest.md) |  | 

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


## ExtendedcommunitylistsPut

> ExtendedcommunitylistsPut(ctx).ChangesetName(changesetName).ExtendedcommunitylistsPutRequest(extendedcommunitylistsPutRequest).Execute()

Create Extended Community List



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
	extendedcommunitylistsPutRequest := *openapiclient.NewExtendedcommunitylistsPutRequest() // ExtendedcommunitylistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ExtendedCommunityListsAPI.ExtendedcommunitylistsPut(context.Background()).ChangesetName(changesetName).ExtendedcommunitylistsPutRequest(extendedcommunitylistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExtendedCommunityListsAPI.ExtendedcommunitylistsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiExtendedcommunitylistsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **extendedcommunitylistsPutRequest** | [**ExtendedcommunitylistsPutRequest**](ExtendedcommunitylistsPutRequest.md) |  | 

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

