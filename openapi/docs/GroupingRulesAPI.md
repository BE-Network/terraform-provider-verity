# \GroupingRulesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GroupingrulesDelete**](GroupingRulesAPI.md#GroupingrulesDelete) | **Delete** /groupingrules | Delete Grouping Rule
[**GroupingrulesGet**](GroupingRulesAPI.md#GroupingrulesGet) | **Get** /groupingrules | Get all Grouping Rules
[**GroupingrulesPatch**](GroupingRulesAPI.md#GroupingrulesPatch) | **Patch** /groupingrules | Update Grouping Rule
[**GroupingrulesPut**](GroupingRulesAPI.md#GroupingrulesPut) | **Put** /groupingrules | Create Grouping Rule



## GroupingrulesDelete

> GroupingrulesDelete(ctx).GroupingRulesName(groupingRulesName).ChangesetName(changesetName).Execute()

Delete Grouping Rule



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
	groupingRulesName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.GroupingRulesAPI.GroupingrulesDelete(context.Background()).GroupingRulesName(groupingRulesName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GroupingRulesAPI.GroupingrulesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGroupingrulesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupingRulesName** | **[]string** |  | 
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


## GroupingrulesGet

> GroupingrulesGet(ctx).GroupingRulesName(groupingRulesName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Grouping Rules



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
	groupingRulesName := "groupingRulesName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.GroupingRulesAPI.GroupingrulesGet(context.Background()).GroupingRulesName(groupingRulesName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GroupingRulesAPI.GroupingrulesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGroupingrulesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupingRulesName** | **string** |  | 
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


## GroupingrulesPatch

> GroupingrulesPatch(ctx).ChangesetName(changesetName).GroupingrulesPutRequest(groupingrulesPutRequest).Execute()

Update Grouping Rule



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
	groupingrulesPutRequest := *openapiclient.NewGroupingrulesPutRequest() // GroupingrulesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.GroupingRulesAPI.GroupingrulesPatch(context.Background()).ChangesetName(changesetName).GroupingrulesPutRequest(groupingrulesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GroupingRulesAPI.GroupingrulesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGroupingrulesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **groupingrulesPutRequest** | [**GroupingrulesPutRequest**](GroupingrulesPutRequest.md) |  | 

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


## GroupingrulesPut

> GroupingrulesPut(ctx).ChangesetName(changesetName).GroupingrulesPutRequest(groupingrulesPutRequest).Execute()

Create Grouping Rule



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
	groupingrulesPutRequest := *openapiclient.NewGroupingrulesPutRequest() // GroupingrulesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.GroupingRulesAPI.GroupingrulesPut(context.Background()).ChangesetName(changesetName).GroupingrulesPutRequest(groupingrulesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `GroupingRulesAPI.GroupingrulesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGroupingrulesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **groupingrulesPutRequest** | [**GroupingrulesPutRequest**](GroupingrulesPutRequest.md) |  | 

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

