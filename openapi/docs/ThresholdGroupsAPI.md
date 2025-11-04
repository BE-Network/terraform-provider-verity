# \ThresholdGroupsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ThresholdgroupsDelete**](ThresholdGroupsAPI.md#ThresholdgroupsDelete) | **Delete** /thresholdgroups | Delete Threshold Group
[**ThresholdgroupsGet**](ThresholdGroupsAPI.md#ThresholdgroupsGet) | **Get** /thresholdgroups | Get all Threshold Groups
[**ThresholdgroupsPatch**](ThresholdGroupsAPI.md#ThresholdgroupsPatch) | **Patch** /thresholdgroups | Update Threshold Group
[**ThresholdgroupsPut**](ThresholdGroupsAPI.md#ThresholdgroupsPut) | **Put** /thresholdgroups | Create Threshold Group



## ThresholdgroupsDelete

> ThresholdgroupsDelete(ctx).ThresholdGroupName(thresholdGroupName).ChangesetName(changesetName).Execute()

Delete Threshold Group



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
	thresholdGroupName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ThresholdGroupsAPI.ThresholdgroupsDelete(context.Background()).ThresholdGroupName(thresholdGroupName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ThresholdGroupsAPI.ThresholdgroupsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiThresholdgroupsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **thresholdGroupName** | **[]string** |  | 
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


## ThresholdgroupsGet

> ThresholdgroupsGet(ctx).ThresholdGroupName(thresholdGroupName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Threshold Groups



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
	thresholdGroupName := "thresholdGroupName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ThresholdGroupsAPI.ThresholdgroupsGet(context.Background()).ThresholdGroupName(thresholdGroupName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ThresholdGroupsAPI.ThresholdgroupsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiThresholdgroupsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **thresholdGroupName** | **string** |  | 
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


## ThresholdgroupsPatch

> ThresholdgroupsPatch(ctx).ChangesetName(changesetName).ThresholdgroupsPutRequest(thresholdgroupsPutRequest).Execute()

Update Threshold Group



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
	thresholdgroupsPutRequest := *openapiclient.NewThresholdgroupsPutRequest() // ThresholdgroupsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ThresholdGroupsAPI.ThresholdgroupsPatch(context.Background()).ChangesetName(changesetName).ThresholdgroupsPutRequest(thresholdgroupsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ThresholdGroupsAPI.ThresholdgroupsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiThresholdgroupsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **thresholdgroupsPutRequest** | [**ThresholdgroupsPutRequest**](ThresholdgroupsPutRequest.md) |  | 

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


## ThresholdgroupsPut

> ThresholdgroupsPut(ctx).ChangesetName(changesetName).ThresholdgroupsPutRequest(thresholdgroupsPutRequest).Execute()

Create Threshold Group



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
	thresholdgroupsPutRequest := *openapiclient.NewThresholdgroupsPutRequest() // ThresholdgroupsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ThresholdGroupsAPI.ThresholdgroupsPut(context.Background()).ChangesetName(changesetName).ThresholdgroupsPutRequest(thresholdgroupsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ThresholdGroupsAPI.ThresholdgroupsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiThresholdgroupsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **thresholdgroupsPutRequest** | [**ThresholdgroupsPutRequest**](ThresholdgroupsPutRequest.md) |  | 

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

