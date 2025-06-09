# \ACLsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AclsDelete**](ACLsAPI.md#AclsDelete) | **Delete** /acls | Delete IP Filter
[**AclsGet**](ACLsAPI.md#AclsGet) | **Get** /acls | Get all IP Filters
[**AclsPatch**](ACLsAPI.md#AclsPatch) | **Patch** /acls | Update IP Filter
[**AclsPut**](ACLsAPI.md#AclsPut) | **Put** /acls | Create IP Filter



## AclsDelete

> AclsDelete(ctx).IpFilterName(ipFilterName).IpVersion(ipVersion).ChangesetName(changesetName).Execute()

Delete IP Filter



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
	ipFilterName := []string{"Inner_example"} // []string | 
	ipVersion := "ipVersion_example" // string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ACLsAPI.AclsDelete(context.Background()).IpFilterName(ipFilterName).IpVersion(ipVersion).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ACLsAPI.AclsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAclsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipFilterName** | **[]string** |  | 
 **ipVersion** | **string** |  | 
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


## AclsGet

> AclsGet(ctx).IpVersion(ipVersion).IpFilterName(ipFilterName).IncludeData(includeData).Execute()

Get all IP Filters



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
	ipVersion := "ipVersion_example" // string | 
	ipFilterName := "ipFilterName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ACLsAPI.AclsGet(context.Background()).IpVersion(ipVersion).IpFilterName(ipFilterName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ACLsAPI.AclsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAclsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipVersion** | **string** |  | 
 **ipFilterName** | **string** |  | 
 **includeData** | **bool** |  | 

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


## AclsPatch

> AclsPatch(ctx).IpVersion(ipVersion).ChangesetName(changesetName).AclsPutRequest(aclsPutRequest).Execute()

Update IP Filter



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
	ipVersion := "ipVersion_example" // string | 
	changesetName := "changesetName_example" // string |  (optional)
	aclsPutRequest := *openapiclient.NewAclsPutRequest() // AclsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ACLsAPI.AclsPatch(context.Background()).IpVersion(ipVersion).ChangesetName(changesetName).AclsPutRequest(aclsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ACLsAPI.AclsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAclsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipVersion** | **string** |  | 
 **changesetName** | **string** |  | 
 **aclsPutRequest** | [**AclsPutRequest**](AclsPutRequest.md) |  | 

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


## AclsPut

> AclsPut(ctx).IpVersion(ipVersion).ChangesetName(changesetName).AclsPutRequest(aclsPutRequest).Execute()

Create IP Filter



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
	ipVersion := "ipVersion_example" // string | 
	changesetName := "changesetName_example" // string |  (optional)
	aclsPutRequest := *openapiclient.NewAclsPutRequest() // AclsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ACLsAPI.AclsPut(context.Background()).IpVersion(ipVersion).ChangesetName(changesetName).AclsPutRequest(aclsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ACLsAPI.AclsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAclsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipVersion** | **string** |  | 
 **changesetName** | **string** |  | 
 **aclsPutRequest** | [**AclsPutRequest**](AclsPutRequest.md) |  | 

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

