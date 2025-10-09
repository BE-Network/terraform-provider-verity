# \PBRoutingAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PolicybasedroutingDelete**](PBRoutingAPI.md#PolicybasedroutingDelete) | **Delete** /policybasedrouting | Delete PB Routing object
[**PolicybasedroutingGet**](PBRoutingAPI.md#PolicybasedroutingGet) | **Get** /policybasedrouting | Get all PB Routing objects
[**PolicybasedroutingPatch**](PBRoutingAPI.md#PolicybasedroutingPatch) | **Patch** /policybasedrouting | Update PB Routing object
[**PolicybasedroutingPut**](PBRoutingAPI.md#PolicybasedroutingPut) | **Put** /policybasedrouting | Create PB Routing object



## PolicybasedroutingDelete

> PolicybasedroutingDelete(ctx).PbRoutingName(pbRoutingName).ChangesetName(changesetName).Execute()

Delete PB Routing object



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
	pbRoutingName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingAPI.PolicybasedroutingDelete(context.Background()).PbRoutingName(pbRoutingName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingAPI.PolicybasedroutingDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pbRoutingName** | **[]string** |  | 
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


## PolicybasedroutingGet

> PolicybasedroutingGet(ctx).PbRoutingName(pbRoutingName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all PB Routing objects



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
	pbRoutingName := "pbRoutingName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingAPI.PolicybasedroutingGet(context.Background()).PbRoutingName(pbRoutingName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingAPI.PolicybasedroutingGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pbRoutingName** | **string** |  | 
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


## PolicybasedroutingPatch

> PolicybasedroutingPatch(ctx).ChangesetName(changesetName).PolicybasedroutingPutRequest(policybasedroutingPutRequest).Execute()

Update PB Routing object



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
	policybasedroutingPutRequest := *openapiclient.NewPolicybasedroutingPutRequest() // PolicybasedroutingPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingAPI.PolicybasedroutingPatch(context.Background()).ChangesetName(changesetName).PolicybasedroutingPutRequest(policybasedroutingPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingAPI.PolicybasedroutingPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **policybasedroutingPutRequest** | [**PolicybasedroutingPutRequest**](PolicybasedroutingPutRequest.md) |  | 

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


## PolicybasedroutingPut

> PolicybasedroutingPut(ctx).ChangesetName(changesetName).PolicybasedroutingPutRequest(policybasedroutingPutRequest).Execute()

Create PB Routing object



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
	policybasedroutingPutRequest := *openapiclient.NewPolicybasedroutingPutRequest() // PolicybasedroutingPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingAPI.PolicybasedroutingPut(context.Background()).ChangesetName(changesetName).PolicybasedroutingPutRequest(policybasedroutingPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingAPI.PolicybasedroutingPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **policybasedroutingPutRequest** | [**PolicybasedroutingPutRequest**](PolicybasedroutingPutRequest.md) |  | 

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

