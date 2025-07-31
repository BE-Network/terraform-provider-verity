# \DeviceControllersAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DevicecontrollersDelete**](DeviceControllersAPI.md#DevicecontrollersDelete) | **Delete** /devicecontrollers | Delete Device Controllers
[**DevicecontrollersGet**](DeviceControllersAPI.md#DevicecontrollersGet) | **Get** /devicecontrollers | Get all Device Controllers
[**DevicecontrollersPatch**](DeviceControllersAPI.md#DevicecontrollersPatch) | **Patch** /devicecontrollers | Update Device Controller
[**DevicecontrollersPut**](DeviceControllersAPI.md#DevicecontrollersPut) | **Put** /devicecontrollers | Create Device Controller



## DevicecontrollersDelete

> DevicecontrollersDelete(ctx).DeviceControllerName(deviceControllerName).ChangesetName(changesetName).Execute()

Delete Device Controllers



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
	deviceControllerName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceControllersAPI.DevicecontrollersDelete(context.Background()).DeviceControllerName(deviceControllerName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceControllersAPI.DevicecontrollersDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicecontrollersDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceControllerName** | **[]string** |  | 
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


## DevicecontrollersGet

> DevicecontrollersGet(ctx).DeviceControllerName(deviceControllerName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Device Controllers



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
	deviceControllerName := "deviceControllerName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceControllersAPI.DevicecontrollersGet(context.Background()).DeviceControllerName(deviceControllerName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceControllersAPI.DevicecontrollersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicecontrollersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceControllerName** | **string** |  | 
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


## DevicecontrollersPatch

> DevicecontrollersPatch(ctx).ChangesetName(changesetName).DevicecontrollersPutRequest(devicecontrollersPutRequest).Execute()

Update Device Controller



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
	devicecontrollersPutRequest := *openapiclient.NewDevicecontrollersPutRequest() // DevicecontrollersPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceControllersAPI.DevicecontrollersPatch(context.Background()).ChangesetName(changesetName).DevicecontrollersPutRequest(devicecontrollersPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceControllersAPI.DevicecontrollersPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicecontrollersPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **devicecontrollersPutRequest** | [**DevicecontrollersPutRequest**](DevicecontrollersPutRequest.md) |  | 

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


## DevicecontrollersPut

> DevicecontrollersPut(ctx).ChangesetName(changesetName).DevicecontrollersPutRequest(devicecontrollersPutRequest).Execute()

Create Device Controller



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
	devicecontrollersPutRequest := *openapiclient.NewDevicecontrollersPutRequest() // DevicecontrollersPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceControllersAPI.DevicecontrollersPut(context.Background()).ChangesetName(changesetName).DevicecontrollersPutRequest(devicecontrollersPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceControllersAPI.DevicecontrollersPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicecontrollersPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **devicecontrollersPutRequest** | [**DevicecontrollersPutRequest**](DevicecontrollersPutRequest.md) |  | 

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

