# \IPv6PrefixListsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ipv6prefixlistsDelete**](IPv6PrefixListsAPI.md#Ipv6prefixlistsDelete) | **Delete** /ipv6prefixlists | Delete IPv6 Prefix List
[**Ipv6prefixlistsGet**](IPv6PrefixListsAPI.md#Ipv6prefixlistsGet) | **Get** /ipv6prefixlists | Get all IPv6 Prefix Lists
[**Ipv6prefixlistsPatch**](IPv6PrefixListsAPI.md#Ipv6prefixlistsPatch) | **Patch** /ipv6prefixlists | Update IPv6 Prefix List
[**Ipv6prefixlistsPut**](IPv6PrefixListsAPI.md#Ipv6prefixlistsPut) | **Put** /ipv6prefixlists | Create IPv6 Prefix List



## Ipv6prefixlistsDelete

> Ipv6prefixlistsDelete(ctx).Ipv6PrefixListName(ipv6PrefixListName).ChangesetName(changesetName).Execute()

Delete IPv6 Prefix List



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
	ipv6PrefixListName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6PrefixListsAPI.Ipv6prefixlistsDelete(context.Background()).Ipv6PrefixListName(ipv6PrefixListName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6PrefixListsAPI.Ipv6prefixlistsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6prefixlistsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv6PrefixListName** | **[]string** |  | 
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


## Ipv6prefixlistsGet

> Ipv6prefixlistsGet(ctx).Ipv6PrefixListName(ipv6PrefixListName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all IPv6 Prefix Lists



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
	ipv6PrefixListName := "ipv6PrefixListName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6PrefixListsAPI.Ipv6prefixlistsGet(context.Background()).Ipv6PrefixListName(ipv6PrefixListName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6PrefixListsAPI.Ipv6prefixlistsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6prefixlistsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv6PrefixListName** | **string** |  | 
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


## Ipv6prefixlistsPatch

> Ipv6prefixlistsPatch(ctx).ChangesetName(changesetName).Ipv6prefixlistsPutRequest(ipv6prefixlistsPutRequest).Execute()

Update IPv6 Prefix List



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
	ipv6prefixlistsPutRequest := *openapiclient.NewIpv6prefixlistsPutRequest() // Ipv6prefixlistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6PrefixListsAPI.Ipv6prefixlistsPatch(context.Background()).ChangesetName(changesetName).Ipv6prefixlistsPutRequest(ipv6prefixlistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6PrefixListsAPI.Ipv6prefixlistsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6prefixlistsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv6prefixlistsPutRequest** | [**Ipv6prefixlistsPutRequest**](Ipv6prefixlistsPutRequest.md) |  | 

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


## Ipv6prefixlistsPut

> Ipv6prefixlistsPut(ctx).ChangesetName(changesetName).Ipv6prefixlistsPutRequest(ipv6prefixlistsPutRequest).Execute()

Create IPv6 Prefix List



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
	ipv6prefixlistsPutRequest := *openapiclient.NewIpv6prefixlistsPutRequest() // Ipv6prefixlistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv6PrefixListsAPI.Ipv6prefixlistsPut(context.Background()).ChangesetName(changesetName).Ipv6prefixlistsPutRequest(ipv6prefixlistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv6PrefixListsAPI.Ipv6prefixlistsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv6prefixlistsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv6prefixlistsPutRequest** | [**Ipv6prefixlistsPutRequest**](Ipv6prefixlistsPutRequest.md) |  | 

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

