# \ReadModeAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ReadmodeGet**](ReadModeAPI.md#ReadmodeGet) | **Get** /readmode | Get list of read-only switchpoints or pods
[**ReadmodePut**](ReadModeAPI.md#ReadmodePut) | **Put** /readmode | Change read-only mode of switchpoint or pod



## ReadmodeGet

> ReadmodeGet(ctx).Selector(selector).ReadOnlyMode(readOnlyMode).Execute()

Get list of read-only switchpoints or pods



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
	selector := "selector_example" // string | 
	readOnlyMode := true // bool | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ReadModeAPI.ReadmodeGet(context.Background()).Selector(selector).ReadOnlyMode(readOnlyMode).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReadModeAPI.ReadmodeGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReadmodeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **selector** | **string** |  | 
 **readOnlyMode** | **bool** |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ReadmodePut

> ReadmodePut(ctx).ReadOnlyMode(readOnlyMode).Pod(pod).Switchpoint(switchpoint).Execute()

Change read-only mode of switchpoint or pod



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
	readOnlyMode := true // bool | 
	pod := []string{"Inner_example"} // []string |  (optional)
	switchpoint := []string{"Inner_example"} // []string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ReadModeAPI.ReadmodePut(context.Background()).ReadOnlyMode(readOnlyMode).Pod(pod).Switchpoint(switchpoint).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ReadModeAPI.ReadmodePut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiReadmodePutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **readOnlyMode** | **bool** |  | 
 **pod** | **[]string** |  | 
 **switchpoint** | **[]string** |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

