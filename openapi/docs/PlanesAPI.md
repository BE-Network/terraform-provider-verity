# \PlanesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PlanesDelete**](PlanesAPI.md#PlanesDelete) | **Delete** /planes | Delete plane
[**PlanesGet**](PlanesAPI.md#PlanesGet) | **Get** /planes | Get all planes
[**PlanesPatch**](PlanesAPI.md#PlanesPatch) | **Patch** /planes | Update plane
[**PlanesPut**](PlanesAPI.md#PlanesPut) | **Put** /planes | Create plane



## PlanesDelete

> PlanesDelete(ctx).PlaneName(planeName).ChangesetName(changesetName).Execute()

Delete plane



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
	planeName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlanesAPI.PlanesDelete(context.Background()).PlaneName(planeName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlanesAPI.PlanesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlanesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **planeName** | **[]string** |  | 
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


## PlanesGet

> PlanesGet(ctx).PlaneName(planeName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all planes



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
	planeName := "planeName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlanesAPI.PlanesGet(context.Background()).PlaneName(planeName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlanesAPI.PlanesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlanesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **planeName** | **string** |  | 
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


## PlanesPatch

> PlanesPatch(ctx).ChangesetName(changesetName).PlanesPutRequest(planesPutRequest).Execute()

Update plane



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
	planesPutRequest := *openapiclient.NewPlanesPutRequest() // PlanesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlanesAPI.PlanesPatch(context.Background()).ChangesetName(changesetName).PlanesPutRequest(planesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlanesAPI.PlanesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlanesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **planesPutRequest** | [**PlanesPutRequest**](PlanesPutRequest.md) |  | 

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


## PlanesPut

> PlanesPut(ctx).ChangesetName(changesetName).PlanesPutRequest(planesPutRequest).Execute()

Create plane



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
	planesPutRequest := *openapiclient.NewPlanesPutRequest() // PlanesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PlanesAPI.PlanesPut(context.Background()).ChangesetName(changesetName).PlanesPutRequest(planesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PlanesAPI.PlanesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPlanesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **planesPutRequest** | [**PlanesPutRequest**](PlanesPutRequest.md) |  | 

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

