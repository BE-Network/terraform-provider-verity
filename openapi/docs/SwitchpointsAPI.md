# \SwitchpointsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SwitchpointsCurrentconfigGet**](SwitchpointsAPI.md#SwitchpointsCurrentconfigGet) | **Get** /switchpoints/currentconfig | Get all Switchpoint current configs



## SwitchpointsCurrentconfigGet

> SwitchpointsCurrentconfigGet(ctx).SwitchpointName(switchpointName).Execute()

Get all Switchpoint current configs



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
	switchpointName := "switchpointName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsCurrentconfigGet(context.Background()).SwitchpointName(switchpointName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsCurrentconfigGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsCurrentconfigGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchpointName** | **string** |  | 

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

