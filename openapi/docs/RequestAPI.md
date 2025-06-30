# \RequestAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RequestGet**](RequestAPI.md#RequestGet) | **Get** /request | Get miscellaneous info about the system



## RequestGet

> RequestGet(ctx).Query(query).Setname(setname).Site(site).DevId(devId).EndpointId(endpointId).EndpointName(endpointName).Start(start).Stop(stop).Endpoint(endpoint).Execute()

Get miscellaneous info about the system



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
	query := "query_example" // string | 
	setname := "setname_example" // string |  (optional)
	site := "site_example" // string |  (optional)
	devId := int32(56) // int32 |  (optional)
	endpointId := int32(56) // int32 |  (optional)
	endpointName := "endpointName_example" // string |  (optional)
	start := "start_example" // string |  (optional)
	stop := "stop_example" // string |  (optional)
	endpoint := "endpoint_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RequestAPI.RequestGet(context.Background()).Query(query).Setname(setname).Site(site).DevId(devId).EndpointId(endpointId).EndpointName(endpointName).Start(start).Stop(stop).Endpoint(endpoint).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RequestAPI.RequestGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRequestGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **query** | **string** |  | 
 **setname** | **string** |  | 
 **site** | **string** |  | 
 **devId** | **int32** |  | 
 **endpointId** | **int32** |  | 
 **endpointName** | **string** |  | 
 **start** | **string** |  | 
 **stop** | **string** |  | 
 **endpoint** | **string** |  | 

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

