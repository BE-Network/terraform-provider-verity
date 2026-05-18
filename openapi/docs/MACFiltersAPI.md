# \MACFiltersAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**MacfiltersDelete**](MACFiltersAPI.md#MacfiltersDelete) | **Delete** /macfilters | Delete MAC Filter
[**MacfiltersGet**](MACFiltersAPI.md#MacfiltersGet) | **Get** /macfilters | Get all MAC Filters
[**MacfiltersPatch**](MACFiltersAPI.md#MacfiltersPatch) | **Patch** /macfilters | Update MAC Filter
[**MacfiltersPut**](MACFiltersAPI.md#MacfiltersPut) | **Put** /macfilters | Create MAC Filter



## MacfiltersDelete

> MacfiltersDelete(ctx).MacFilterName(macFilterName).ChangesetName(changesetName).Execute()

Delete MAC Filter



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
	macFilterName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MACFiltersAPI.MacfiltersDelete(context.Background()).MacFilterName(macFilterName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MACFiltersAPI.MacfiltersDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMacfiltersDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **macFilterName** | **[]string** |  | 
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


## MacfiltersGet

> MacfiltersGet(ctx).MacFilterName(macFilterName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all MAC Filters



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
	macFilterName := "macFilterName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MACFiltersAPI.MacfiltersGet(context.Background()).MacFilterName(macFilterName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MACFiltersAPI.MacfiltersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMacfiltersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **macFilterName** | **string** |  | 
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


## MacfiltersPatch

> MacfiltersPatch(ctx).ChangesetName(changesetName).MacfiltersPutRequest(macfiltersPutRequest).Execute()

Update MAC Filter



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
	macfiltersPutRequest := *openapiclient.NewMacfiltersPutRequest() // MacfiltersPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MACFiltersAPI.MacfiltersPatch(context.Background()).ChangesetName(changesetName).MacfiltersPutRequest(macfiltersPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MACFiltersAPI.MacfiltersPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMacfiltersPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **macfiltersPutRequest** | [**MacfiltersPutRequest**](MacfiltersPutRequest.md) |  | 

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


## MacfiltersPut

> MacfiltersPut(ctx).ChangesetName(changesetName).MacfiltersPutRequest(macfiltersPutRequest).Execute()

Create MAC Filter



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
	macfiltersPutRequest := *openapiclient.NewMacfiltersPutRequest() // MacfiltersPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MACFiltersAPI.MacfiltersPut(context.Background()).ChangesetName(changesetName).MacfiltersPutRequest(macfiltersPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MACFiltersAPI.MacfiltersPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMacfiltersPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **macfiltersPutRequest** | [**MacfiltersPutRequest**](MacfiltersPutRequest.md) |  | 

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

