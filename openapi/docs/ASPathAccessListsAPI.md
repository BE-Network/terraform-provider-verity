# \ASPathAccessListsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AspathaccesslistsDelete**](ASPathAccessListsAPI.md#AspathaccesslistsDelete) | **Delete** /aspathaccesslists | Delete AS Path Access Lists
[**AspathaccesslistsGet**](ASPathAccessListsAPI.md#AspathaccesslistsGet) | **Get** /aspathaccesslists | Get all AS Path Access Lists
[**AspathaccesslistsPatch**](ASPathAccessListsAPI.md#AspathaccesslistsPatch) | **Patch** /aspathaccesslists | Update AS Path Access List
[**AspathaccesslistsPut**](ASPathAccessListsAPI.md#AspathaccesslistsPut) | **Put** /aspathaccesslists | Create AS Path Access List



## AspathaccesslistsDelete

> AspathaccesslistsDelete(ctx).AsPathAccessListName(asPathAccessListName).ChangesetName(changesetName).Execute()

Delete AS Path Access Lists



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
	asPathAccessListName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ASPathAccessListsAPI.AspathaccesslistsDelete(context.Background()).AsPathAccessListName(asPathAccessListName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ASPathAccessListsAPI.AspathaccesslistsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAspathaccesslistsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **asPathAccessListName** | **[]string** |  | 
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


## AspathaccesslistsGet

> AspathaccesslistsGet(ctx).AsPathAccessListName(asPathAccessListName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all AS Path Access Lists



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
	asPathAccessListName := "asPathAccessListName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ASPathAccessListsAPI.AspathaccesslistsGet(context.Background()).AsPathAccessListName(asPathAccessListName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ASPathAccessListsAPI.AspathaccesslistsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAspathaccesslistsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **asPathAccessListName** | **string** |  | 
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


## AspathaccesslistsPatch

> AspathaccesslistsPatch(ctx).ChangesetName(changesetName).AspathaccesslistsPutRequest(aspathaccesslistsPutRequest).Execute()

Update AS Path Access List



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
	aspathaccesslistsPutRequest := *openapiclient.NewAspathaccesslistsPutRequest() // AspathaccesslistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ASPathAccessListsAPI.AspathaccesslistsPatch(context.Background()).ChangesetName(changesetName).AspathaccesslistsPutRequest(aspathaccesslistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ASPathAccessListsAPI.AspathaccesslistsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAspathaccesslistsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **aspathaccesslistsPutRequest** | [**AspathaccesslistsPutRequest**](AspathaccesslistsPutRequest.md) |  | 

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


## AspathaccesslistsPut

> AspathaccesslistsPut(ctx).ChangesetName(changesetName).AspathaccesslistsPutRequest(aspathaccesslistsPutRequest).Execute()

Create AS Path Access List



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
	aspathaccesslistsPutRequest := *openapiclient.NewAspathaccesslistsPutRequest() // AspathaccesslistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ASPathAccessListsAPI.AspathaccesslistsPut(context.Background()).ChangesetName(changesetName).AspathaccesslistsPutRequest(aspathaccesslistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ASPathAccessListsAPI.AspathaccesslistsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAspathaccesslistsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **aspathaccesslistsPutRequest** | [**AspathaccesslistsPutRequest**](AspathaccesslistsPutRequest.md) |  | 

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

