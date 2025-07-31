# \IPv4PrefixListsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Ipv4prefixlistsDelete**](IPv4PrefixListsAPI.md#Ipv4prefixlistsDelete) | **Delete** /ipv4prefixlists | Delete IPv4 Prefix List
[**Ipv4prefixlistsGet**](IPv4PrefixListsAPI.md#Ipv4prefixlistsGet) | **Get** /ipv4prefixlists | Get all IPv4 Prefix Lists
[**Ipv4prefixlistsPatch**](IPv4PrefixListsAPI.md#Ipv4prefixlistsPatch) | **Patch** /ipv4prefixlists | Update IPv4 Prefix List
[**Ipv4prefixlistsPut**](IPv4PrefixListsAPI.md#Ipv4prefixlistsPut) | **Put** /ipv4prefixlists | Create IPv4 Prefix List



## Ipv4prefixlistsDelete

> Ipv4prefixlistsDelete(ctx).Ipv4PrefixListName(ipv4PrefixListName).ChangesetName(changesetName).Execute()

Delete IPv4 Prefix List



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
	ipv4PrefixListName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4PrefixListsAPI.Ipv4prefixlistsDelete(context.Background()).Ipv4PrefixListName(ipv4PrefixListName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4PrefixListsAPI.Ipv4prefixlistsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4prefixlistsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv4PrefixListName** | **[]string** |  | 
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


## Ipv4prefixlistsGet

> Ipv4prefixlistsGet(ctx).Ipv4PrefixListName(ipv4PrefixListName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all IPv4 Prefix Lists



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
	ipv4PrefixListName := "ipv4PrefixListName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4PrefixListsAPI.Ipv4prefixlistsGet(context.Background()).Ipv4PrefixListName(ipv4PrefixListName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4PrefixListsAPI.Ipv4prefixlistsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4prefixlistsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ipv4PrefixListName** | **string** |  | 
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


## Ipv4prefixlistsPatch

> Ipv4prefixlistsPatch(ctx).ChangesetName(changesetName).Ipv4prefixlistsPutRequest(ipv4prefixlistsPutRequest).Execute()

Update IPv4 Prefix List



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
	ipv4prefixlistsPutRequest := *openapiclient.NewIpv4prefixlistsPutRequest() // Ipv4prefixlistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4PrefixListsAPI.Ipv4prefixlistsPatch(context.Background()).ChangesetName(changesetName).Ipv4prefixlistsPutRequest(ipv4prefixlistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4PrefixListsAPI.Ipv4prefixlistsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4prefixlistsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv4prefixlistsPutRequest** | [**Ipv4prefixlistsPutRequest**](Ipv4prefixlistsPutRequest.md) |  | 

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


## Ipv4prefixlistsPut

> Ipv4prefixlistsPut(ctx).ChangesetName(changesetName).Ipv4prefixlistsPutRequest(ipv4prefixlistsPutRequest).Execute()

Create IPv4 Prefix List



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
	ipv4prefixlistsPutRequest := *openapiclient.NewIpv4prefixlistsPutRequest() // Ipv4prefixlistsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.IPv4PrefixListsAPI.Ipv4prefixlistsPut(context.Background()).ChangesetName(changesetName).Ipv4prefixlistsPutRequest(ipv4prefixlistsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `IPv4PrefixListsAPI.Ipv4prefixlistsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIpv4prefixlistsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ipv4prefixlistsPutRequest** | [**Ipv4prefixlistsPutRequest**](Ipv4prefixlistsPutRequest.md) |  | 

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

