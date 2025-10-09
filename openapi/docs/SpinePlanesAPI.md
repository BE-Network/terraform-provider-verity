# \SpinePlanesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SpineplanesDelete**](SpinePlanesAPI.md#SpineplanesDelete) | **Delete** /spineplanes | Delete Spine Plane
[**SpineplanesGet**](SpinePlanesAPI.md#SpineplanesGet) | **Get** /spineplanes | Get all Spine Planes
[**SpineplanesPatch**](SpinePlanesAPI.md#SpineplanesPatch) | **Patch** /spineplanes | Update Spine Plane
[**SpineplanesPut**](SpinePlanesAPI.md#SpineplanesPut) | **Put** /spineplanes | Create Spine Plane



## SpineplanesDelete

> SpineplanesDelete(ctx).SpinePlaneName(spinePlaneName).ChangesetName(changesetName).Execute()

Delete Spine Plane



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
	spinePlaneName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SpinePlanesAPI.SpineplanesDelete(context.Background()).SpinePlaneName(spinePlaneName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SpinePlanesAPI.SpineplanesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSpineplanesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **spinePlaneName** | **[]string** |  | 
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


## SpineplanesGet

> SpineplanesGet(ctx).SpinePlaneName(spinePlaneName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Spine Planes



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
	spinePlaneName := "spinePlaneName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SpinePlanesAPI.SpineplanesGet(context.Background()).SpinePlaneName(spinePlaneName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SpinePlanesAPI.SpineplanesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSpineplanesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **spinePlaneName** | **string** |  | 
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


## SpineplanesPatch

> SpineplanesPatch(ctx).ChangesetName(changesetName).SpineplanesPutRequest(spineplanesPutRequest).Execute()

Update Spine Plane



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
	spineplanesPutRequest := *openapiclient.NewSpineplanesPutRequest() // SpineplanesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SpinePlanesAPI.SpineplanesPatch(context.Background()).ChangesetName(changesetName).SpineplanesPutRequest(spineplanesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SpinePlanesAPI.SpineplanesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSpineplanesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **spineplanesPutRequest** | [**SpineplanesPutRequest**](SpineplanesPutRequest.md) |  | 

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


## SpineplanesPut

> SpineplanesPut(ctx).ChangesetName(changesetName).SpineplanesPutRequest(spineplanesPutRequest).Execute()

Create Spine Plane



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
	spineplanesPutRequest := *openapiclient.NewSpineplanesPutRequest() // SpineplanesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SpinePlanesAPI.SpineplanesPut(context.Background()).ChangesetName(changesetName).SpineplanesPutRequest(spineplanesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SpinePlanesAPI.SpineplanesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSpineplanesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **spineplanesPutRequest** | [**SpineplanesPutRequest**](SpineplanesPutRequest.md) |  | 

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

