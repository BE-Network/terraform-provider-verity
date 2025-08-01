# \IPv4ListFiltersAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ipv4listsDelete**](IPv4ListFiltersAPI.md#Ipv4listsDelete) | **Delete** /ipv4lists | Delete IPv4 List Filter
[**Ipv4listsGet**](IPv4ListFiltersAPI.md#Ipv4listsGet) | **Get** /ipv4lists | Get all IPv4 List Filters
[**Ipv4listsPatch**](IPv4ListFiltersAPI.md#Ipv4listsPatch) | **Patch** /ipv4lists | Update IPv4 List Filter
[**Ipv4listsPut**](IPv4ListFiltersAPI.md#Ipv4listsPut) | **Put** /ipv4lists | Create IPv4 List Filter



## Ipv4listsDelete

> Ipv4listsDelete(ctx).Ipv4ListFilterName(ipv4ListFilterName).ChangesetName(changesetName).Execute()

Delete IPv4 List Filter



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
	ipv4ListFilterName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4ListFiltersAPI.Ipv4listsDelete(context.Background()).Ipv4ListFilterName(ipv4ListFilterName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4ListFiltersAPI.Ipv4listsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4listsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv4ListFilterName** | **[]string** |  | 
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


## Ipv4listsGet

> Ipv4listsGet(ctx).Ipv4ListFilterName(ipv4ListFilterName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all IPv4 List Filters



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
	ipv4ListFilterName := "ipv4ListFilterName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4ListFiltersAPI.Ipv4listsGet(context.Background()).Ipv4ListFilterName(ipv4ListFilterName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4ListFiltersAPI.Ipv4listsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4listsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv4ListFilterName** | **string** |  | 
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


## Ipv4listsPatch

> Ipv4listsPatch(ctx).ChangesetName(changesetName).Ipv4listsPutRequest(ipv4listsPutRequest).Execute()

Update IPv4 List Filter



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
	ipv4listsPutRequest := *openapiclient.NewIpv4listsPutRequest() // Ipv4listsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4ListFiltersAPI.Ipv4listsPatch(context.Background()).ChangesetName(changesetName).Ipv4listsPutRequest(ipv4listsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4ListFiltersAPI.Ipv4listsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4listsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv4listsPutRequest** | [**Ipv4listsPutRequest**](Ipv4listsPutRequest.md) |  | 

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


## Ipv4listsPut

> Ipv4listsPut(ctx).ChangesetName(changesetName).Ipv4listsPutRequest(ipv4listsPutRequest).Execute()

Create IPv4 List Filter



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
	ipv4listsPutRequest := *openapiclient.NewIpv4listsPutRequest() // Ipv4listsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4ListFiltersAPI.Ipv4listsPut(context.Background()).ChangesetName(changesetName).Ipv4listsPutRequest(ipv4listsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4ListFiltersAPI.Ipv4listsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4listsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv4listsPutRequest** | [**Ipv4listsPutRequest**](Ipv4listsPutRequest.md) |  | 

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

