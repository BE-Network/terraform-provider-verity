# \SuperSpineGroupsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SspgroupsDelete**](SuperSpineGroupsAPI.md#SspgroupsDelete) | **Delete** /sspgroups | Delete SuperSpine Group
[**SspgroupsGet**](SuperSpineGroupsAPI.md#SspgroupsGet) | **Get** /sspgroups | Get all SuperSpine Groups
[**SspgroupsPatch**](SuperSpineGroupsAPI.md#SspgroupsPatch) | **Patch** /sspgroups | Update SuperSpine Group
[**SspgroupsPut**](SuperSpineGroupsAPI.md#SspgroupsPut) | **Put** /sspgroups | Create SuperSpine Group



## SspgroupsDelete

> SspgroupsDelete(ctx).SuperspineGroupName(superspineGroupName).ChangesetName(changesetName).Execute()

Delete SuperSpine Group



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
	superspineGroupName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SuperSpineGroupsAPI.SspgroupsDelete(context.Background()).SuperspineGroupName(superspineGroupName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SuperSpineGroupsAPI.SspgroupsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSspgroupsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **superspineGroupName** | **[]string** |  | 
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


## SspgroupsGet

> SspgroupsGet(ctx).SuperspineGroupName(superspineGroupName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all SuperSpine Groups



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
	superspineGroupName := "superspineGroupName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SuperSpineGroupsAPI.SspgroupsGet(context.Background()).SuperspineGroupName(superspineGroupName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SuperSpineGroupsAPI.SspgroupsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSspgroupsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **superspineGroupName** | **string** |  | 
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


## SspgroupsPatch

> SspgroupsPatch(ctx).ChangesetName(changesetName).SspgroupsPutRequest(sspgroupsPutRequest).Execute()

Update SuperSpine Group



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
	sspgroupsPutRequest := *openapiclient.NewSspgroupsPutRequest() // SspgroupsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SuperSpineGroupsAPI.SspgroupsPatch(context.Background()).ChangesetName(changesetName).SspgroupsPutRequest(sspgroupsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SuperSpineGroupsAPI.SspgroupsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSspgroupsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **sspgroupsPutRequest** | [**SspgroupsPutRequest**](SspgroupsPutRequest.md) |  | 

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


## SspgroupsPut

> SspgroupsPut(ctx).ChangesetName(changesetName).SspgroupsPutRequest(sspgroupsPutRequest).Execute()

Create SuperSpine Group



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
	sspgroupsPutRequest := *openapiclient.NewSspgroupsPutRequest() // SspgroupsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SuperSpineGroupsAPI.SspgroupsPut(context.Background()).ChangesetName(changesetName).SspgroupsPutRequest(sspgroupsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SuperSpineGroupsAPI.SspgroupsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSspgroupsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **sspgroupsPutRequest** | [**SspgroupsPutRequest**](SspgroupsPutRequest.md) |  | 

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

