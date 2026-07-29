# \FabricsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**FabricsDelete**](FabricsAPI.md#FabricsDelete) | **Delete** /fabrics | Delete Fabric
[**FabricsGet**](FabricsAPI.md#FabricsGet) | **Get** /fabrics | Get all Fabrics
[**FabricsPatch**](FabricsAPI.md#FabricsPatch) | **Patch** /fabrics | Update Fabric
[**FabricsPut**](FabricsAPI.md#FabricsPut) | **Put** /fabrics | Create Fabric



## FabricsDelete

> FabricsDelete(ctx).FabricName(fabricName).ChangesetName(changesetName).Execute()

Delete Fabric



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
	fabricName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FabricsAPI.FabricsDelete(context.Background()).FabricName(fabricName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FabricsAPI.FabricsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiFabricsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fabricName** | **[]string** |  | 
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


## FabricsGet

> FabricsGet(ctx).FabricName(fabricName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Fabrics



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
	fabricName := "fabricName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FabricsAPI.FabricsGet(context.Background()).FabricName(fabricName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FabricsAPI.FabricsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiFabricsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fabricName** | **string** |  | 
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


## FabricsPatch

> FabricsPatch(ctx).ChangesetName(changesetName).FabricsPutRequest(fabricsPutRequest).Execute()

Update Fabric



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
	fabricsPutRequest := *openapiclient.NewFabricsPutRequest() // FabricsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FabricsAPI.FabricsPatch(context.Background()).ChangesetName(changesetName).FabricsPutRequest(fabricsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FabricsAPI.FabricsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiFabricsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **fabricsPutRequest** | [**FabricsPutRequest**](FabricsPutRequest.md) |  | 

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


## FabricsPut

> FabricsPut(ctx).ChangesetName(changesetName).FabricsPutRequest(fabricsPutRequest).Execute()

Create Fabric



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
	fabricsPutRequest := *openapiclient.NewFabricsPutRequest() // FabricsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.FabricsAPI.FabricsPut(context.Background()).ChangesetName(changesetName).FabricsPutRequest(fabricsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `FabricsAPI.FabricsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiFabricsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **fabricsPutRequest** | [**FabricsPutRequest**](FabricsPutRequest.md) |  | 

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

