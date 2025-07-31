# \DeviceSettingsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DevicesettingsDelete**](DeviceSettingsAPI.md#DevicesettingsDelete) | **Delete** /devicesettings | Delete Device Settings
[**DevicesettingsGet**](DeviceSettingsAPI.md#DevicesettingsGet) | **Get** /devicesettings | Get all Device Settings
[**DevicesettingsPatch**](DeviceSettingsAPI.md#DevicesettingsPatch) | **Patch** /devicesettings | Update Device Settings
[**DevicesettingsPut**](DeviceSettingsAPI.md#DevicesettingsPut) | **Put** /devicesettings | Create Device Settings



## DevicesettingsDelete

> DevicesettingsDelete(ctx).EthDeviceProfilesName(ethDeviceProfilesName).ChangesetName(changesetName).Execute()

Delete Device Settings



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
	ethDeviceProfilesName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceSettingsAPI.DevicesettingsDelete(context.Background()).EthDeviceProfilesName(ethDeviceProfilesName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceSettingsAPI.DevicesettingsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicesettingsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ethDeviceProfilesName** | **[]string** |  | 
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


## DevicesettingsGet

> DevicesettingsGet(ctx).EthDeviceProfilesName(ethDeviceProfilesName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Device Settings



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
	ethDeviceProfilesName := "ethDeviceProfilesName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceSettingsAPI.DevicesettingsGet(context.Background()).EthDeviceProfilesName(ethDeviceProfilesName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceSettingsAPI.DevicesettingsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicesettingsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ethDeviceProfilesName** | **string** |  | 
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


## DevicesettingsPatch

> DevicesettingsPatch(ctx).ChangesetName(changesetName).DevicesettingsPutRequest(devicesettingsPutRequest).Execute()

Update Device Settings



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
	devicesettingsPutRequest := *openapiclient.NewDevicesettingsPutRequest() // DevicesettingsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceSettingsAPI.DevicesettingsPatch(context.Background()).ChangesetName(changesetName).DevicesettingsPutRequest(devicesettingsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceSettingsAPI.DevicesettingsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicesettingsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **devicesettingsPutRequest** | [**DevicesettingsPutRequest**](DevicesettingsPutRequest.md) |  | 

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


## DevicesettingsPut

> DevicesettingsPut(ctx).ChangesetName(changesetName).DevicesettingsPutRequest(devicesettingsPutRequest).Execute()

Create Device Settings



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
	devicesettingsPutRequest := *openapiclient.NewDevicesettingsPutRequest() // DevicesettingsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceSettingsAPI.DevicesettingsPut(context.Background()).ChangesetName(changesetName).DevicesettingsPutRequest(devicesettingsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceSettingsAPI.DevicesettingsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicesettingsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **devicesettingsPutRequest** | [**DevicesettingsPutRequest**](DevicesettingsPutRequest.md) |  | 

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

