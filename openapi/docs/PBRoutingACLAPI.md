# \PBRoutingACLAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PolicybasedroutingaclDelete**](PBRoutingACLAPI.md#PolicybasedroutingaclDelete) | **Delete** /policybasedroutingacl | Delete PB Routing ACL
[**PolicybasedroutingaclGet**](PBRoutingACLAPI.md#PolicybasedroutingaclGet) | **Get** /policybasedroutingacl | Get all PB Routing ACLs
[**PolicybasedroutingaclPatch**](PBRoutingACLAPI.md#PolicybasedroutingaclPatch) | **Patch** /policybasedroutingacl | Update PB Routing ACL
[**PolicybasedroutingaclPut**](PBRoutingACLAPI.md#PolicybasedroutingaclPut) | **Put** /policybasedroutingacl | Create PB Routing ACL



## PolicybasedroutingaclDelete

> PolicybasedroutingaclDelete(ctx).PbRoutingAclName(pbRoutingAclName).ChangesetName(changesetName).Execute()

Delete PB Routing ACL



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
	pbRoutingAclName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingACLAPI.PolicybasedroutingaclDelete(context.Background()).PbRoutingAclName(pbRoutingAclName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingACLAPI.PolicybasedroutingaclDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingaclDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pbRoutingAclName** | **[]string** |  | 
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


## PolicybasedroutingaclGet

> PolicybasedroutingaclGet(ctx).PbRoutingAclName(pbRoutingAclName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all PB Routing ACLs



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
	pbRoutingAclName := "pbRoutingAclName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingACLAPI.PolicybasedroutingaclGet(context.Background()).PbRoutingAclName(pbRoutingAclName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingACLAPI.PolicybasedroutingaclGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingaclGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pbRoutingAclName** | **string** |  | 
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


## PolicybasedroutingaclPatch

> PolicybasedroutingaclPatch(ctx).ChangesetName(changesetName).PolicybasedroutingaclPutRequest(policybasedroutingaclPutRequest).Execute()

Update PB Routing ACL



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
	policybasedroutingaclPutRequest := *openapiclient.NewPolicybasedroutingaclPutRequest() // PolicybasedroutingaclPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingACLAPI.PolicybasedroutingaclPatch(context.Background()).ChangesetName(changesetName).PolicybasedroutingaclPutRequest(policybasedroutingaclPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingACLAPI.PolicybasedroutingaclPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingaclPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **policybasedroutingaclPutRequest** | [**PolicybasedroutingaclPutRequest**](PolicybasedroutingaclPutRequest.md) |  | 

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


## PolicybasedroutingaclPut

> PolicybasedroutingaclPut(ctx).ChangesetName(changesetName).PolicybasedroutingaclPutRequest(policybasedroutingaclPutRequest).Execute()

Create PB Routing ACL



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
	policybasedroutingaclPutRequest := *openapiclient.NewPolicybasedroutingaclPutRequest() // PolicybasedroutingaclPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PBRoutingACLAPI.PolicybasedroutingaclPut(context.Background()).ChangesetName(changesetName).PolicybasedroutingaclPutRequest(policybasedroutingaclPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PBRoutingACLAPI.PolicybasedroutingaclPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPolicybasedroutingaclPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **policybasedroutingaclPutRequest** | [**PolicybasedroutingaclPutRequest**](PolicybasedroutingaclPutRequest.md) |  | 

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

