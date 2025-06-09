# \AuthenticatedEthPortsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthenticatedethportsDelete**](AuthenticatedEthPortsAPI.md#AuthenticatedethportsDelete) | **Delete** /authenticatedethports | Delete Authenticated Eth-Port
[**AuthenticatedethportsGet**](AuthenticatedEthPortsAPI.md#AuthenticatedethportsGet) | **Get** /authenticatedethports | Get all Authenticated Eth-Ports
[**AuthenticatedethportsPatch**](AuthenticatedEthPortsAPI.md#AuthenticatedethportsPatch) | **Patch** /authenticatedethports | Update Authenticated Eth-Port
[**AuthenticatedethportsPut**](AuthenticatedEthPortsAPI.md#AuthenticatedethportsPut) | **Put** /authenticatedethports | Create Authenticated Eth-Port



## AuthenticatedethportsDelete

> AuthenticatedethportsDelete(ctx).AuthenticatedEthPortName(authenticatedEthPortName).ChangesetName(changesetName).Execute()

Delete Authenticated Eth-Port



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
	authenticatedEthPortName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AuthenticatedEthPortsAPI.AuthenticatedethportsDelete(context.Background()).AuthenticatedEthPortName(authenticatedEthPortName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticatedEthPortsAPI.AuthenticatedethportsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthenticatedethportsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **authenticatedEthPortName** | **[]string** |  | 
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


## AuthenticatedethportsGet

> AuthenticatedethportsGet(ctx).AuthenticatedEthPortName(authenticatedEthPortName).IncludeData(includeData).Execute()

Get all Authenticated Eth-Ports



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
	authenticatedEthPortName := "authenticatedEthPortName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(context.Background()).AuthenticatedEthPortName(authenticatedEthPortName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticatedEthPortsAPI.AuthenticatedethportsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthenticatedethportsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **authenticatedEthPortName** | **string** |  | 
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


## AuthenticatedethportsPatch

> AuthenticatedethportsPatch(ctx).ChangesetName(changesetName).AuthenticatedethportsPutRequest(authenticatedethportsPutRequest).Execute()

Update Authenticated Eth-Port



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
	authenticatedethportsPutRequest := *openapiclient.NewAuthenticatedethportsPutRequest() // AuthenticatedethportsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AuthenticatedEthPortsAPI.AuthenticatedethportsPatch(context.Background()).ChangesetName(changesetName).AuthenticatedethportsPutRequest(authenticatedethportsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticatedEthPortsAPI.AuthenticatedethportsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthenticatedethportsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **authenticatedethportsPutRequest** | [**AuthenticatedethportsPutRequest**](AuthenticatedethportsPutRequest.md) |  | 

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


## AuthenticatedethportsPut

> AuthenticatedethportsPut(ctx).ChangesetName(changesetName).AuthenticatedethportsPutRequest(authenticatedethportsPutRequest).Execute()

Create Authenticated Eth-Port



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
	authenticatedethportsPutRequest := *openapiclient.NewAuthenticatedethportsPutRequest() // AuthenticatedethportsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.AuthenticatedEthPortsAPI.AuthenticatedethportsPut(context.Background()).ChangesetName(changesetName).AuthenticatedethportsPutRequest(authenticatedethportsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AuthenticatedEthPortsAPI.AuthenticatedethportsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthenticatedethportsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **authenticatedethportsPutRequest** | [**AuthenticatedethportsPutRequest**](AuthenticatedethportsPutRequest.md) |  | 

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

