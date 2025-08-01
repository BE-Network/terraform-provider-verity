# \IPv6ListFiltersAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ipv6listsDelete**](IPv6ListFiltersAPI.md#Ipv6listsDelete) | **Delete** /ipv6lists | Delete IPv6 List Filter
[**Ipv6listsGet**](IPv6ListFiltersAPI.md#Ipv6listsGet) | **Get** /ipv6lists | Get all IPv6 List Filters
[**Ipv6listsPatch**](IPv6ListFiltersAPI.md#Ipv6listsPatch) | **Patch** /ipv6lists | Update IPv6 List Filter
[**Ipv6listsPut**](IPv6ListFiltersAPI.md#Ipv6listsPut) | **Put** /ipv6lists | Create IPv6 List Filter



## Ipv6listsDelete

> Ipv6listsDelete(ctx).Ipv6ListFilterName(ipv6ListFilterName).ChangesetName(changesetName).Execute()

Delete IPv6 List Filter



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
	ipv6ListFilterName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6ListFiltersAPI.Ipv6listsDelete(context.Background()).Ipv6ListFilterName(ipv6ListFilterName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6ListFiltersAPI.Ipv6listsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6listsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv6ListFilterName** | **[]string** |  | 
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


## Ipv6listsGet

> Ipv6listsGet(ctx).Ipv6ListFilterName(ipv6ListFilterName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all IPv6 List Filters



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
	ipv6ListFilterName := "ipv6ListFilterName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6ListFiltersAPI.Ipv6listsGet(context.Background()).Ipv6ListFilterName(ipv6ListFilterName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6ListFiltersAPI.Ipv6listsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6listsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv6ListFilterName** | **string** |  | 
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


## Ipv6listsPatch

> Ipv6listsPatch(ctx).ChangesetName(changesetName).Ipv6listsPutRequest(ipv6listsPutRequest).Execute()

Update IPv6 List Filter



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
	ipv6listsPutRequest := *openapiclient.NewIpv6listsPutRequest() // Ipv6listsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6ListFiltersAPI.Ipv6listsPatch(context.Background()).ChangesetName(changesetName).Ipv6listsPutRequest(ipv6listsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6ListFiltersAPI.Ipv6listsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6listsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv6listsPutRequest** | [**Ipv6listsPutRequest**](Ipv6listsPutRequest.md) |  | 

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


## Ipv6listsPut

> Ipv6listsPut(ctx).ChangesetName(changesetName).Ipv6listsPutRequest(ipv6listsPutRequest).Execute()

Create IPv6 List Filter



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
	ipv6listsPutRequest := *openapiclient.NewIpv6listsPutRequest() // Ipv6listsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6ListFiltersAPI.Ipv6listsPut(context.Background()).ChangesetName(changesetName).Ipv6listsPutRequest(ipv6listsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6ListFiltersAPI.Ipv6listsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6listsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv6listsPutRequest** | [**Ipv6listsPutRequest**](Ipv6listsPutRequest.md) |  | 

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

