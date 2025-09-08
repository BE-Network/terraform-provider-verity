# \SFlowCollectorsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SflowcollectorsDelete**](SFlowCollectorsAPI.md#SflowcollectorsDelete) | **Delete** /sflowcollectors | Delete sFlow Collector
[**SflowcollectorsGet**](SFlowCollectorsAPI.md#SflowcollectorsGet) | **Get** /sflowcollectors | Get all sFlow Collectors
[**SflowcollectorsPatch**](SFlowCollectorsAPI.md#SflowcollectorsPatch) | **Patch** /sflowcollectors | Update sFlow Collector
[**SflowcollectorsPut**](SFlowCollectorsAPI.md#SflowcollectorsPut) | **Put** /sflowcollectors | Create sFlow Collector



## SflowcollectorsDelete

> SflowcollectorsDelete(ctx).SflowCollectorName(sflowCollectorName).ChangesetName(changesetName).Execute()

Delete sFlow Collector



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
	sflowCollectorName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SFlowCollectorsAPI.SflowcollectorsDelete(context.Background()).SflowCollectorName(sflowCollectorName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SFlowCollectorsAPI.SflowcollectorsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSflowcollectorsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sflowCollectorName** | **[]string** |  | 
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


## SflowcollectorsGet

> SflowcollectorsGet(ctx).SflowCollectorName(sflowCollectorName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all sFlow Collectors



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
	sflowCollectorName := "sflowCollectorName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SFlowCollectorsAPI.SflowcollectorsGet(context.Background()).SflowCollectorName(sflowCollectorName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SFlowCollectorsAPI.SflowcollectorsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSflowcollectorsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sflowCollectorName** | **string** |  | 
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


## SflowcollectorsPatch

> SflowcollectorsPatch(ctx).ChangesetName(changesetName).SflowcollectorsPutRequest(sflowcollectorsPutRequest).Execute()

Update sFlow Collector



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
	sflowcollectorsPutRequest := *openapiclient.NewSflowcollectorsPutRequest() // SflowcollectorsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SFlowCollectorsAPI.SflowcollectorsPatch(context.Background()).ChangesetName(changesetName).SflowcollectorsPutRequest(sflowcollectorsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SFlowCollectorsAPI.SflowcollectorsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSflowcollectorsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **sflowcollectorsPutRequest** | [**SflowcollectorsPutRequest**](SflowcollectorsPutRequest.md) |  | 

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


## SflowcollectorsPut

> SflowcollectorsPut(ctx).ChangesetName(changesetName).SflowcollectorsPutRequest(sflowcollectorsPutRequest).Execute()

Create sFlow Collector



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
	sflowcollectorsPutRequest := *openapiclient.NewSflowcollectorsPutRequest() // SflowcollectorsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SFlowCollectorsAPI.SflowcollectorsPut(context.Background()).ChangesetName(changesetName).SflowcollectorsPutRequest(sflowcollectorsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SFlowCollectorsAPI.SflowcollectorsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSflowcollectorsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **sflowcollectorsPutRequest** | [**SflowcollectorsPutRequest**](SflowcollectorsPutRequest.md) |  | 

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

