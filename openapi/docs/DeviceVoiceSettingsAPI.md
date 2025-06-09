# \DeviceVoiceSettingsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DevicevoicesettingsDelete**](DeviceVoiceSettingsAPI.md#DevicevoicesettingsDelete) | **Delete** /devicevoicesettings | Delete tenant
[**DevicevoicesettingsGet**](DeviceVoiceSettingsAPI.md#DevicevoicesettingsGet) | **Get** /devicevoicesettings | Get all Device Voice Settings
[**DevicevoicesettingsPatch**](DeviceVoiceSettingsAPI.md#DevicevoicesettingsPatch) | **Patch** /devicevoicesettings | Update Device Voice Setting
[**DevicevoicesettingsPut**](DeviceVoiceSettingsAPI.md#DevicevoicesettingsPut) | **Put** /devicevoicesettings | Create Device Voice Setting



## DevicevoicesettingsDelete

> DevicevoicesettingsDelete(ctx).DeviceVoiceSettingsName(deviceVoiceSettingsName).ChangesetName(changesetName).Execute()

Delete tenant



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
	deviceVoiceSettingsName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceVoiceSettingsAPI.DevicevoicesettingsDelete(context.Background()).DeviceVoiceSettingsName(deviceVoiceSettingsName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceVoiceSettingsAPI.DevicevoicesettingsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicevoicesettingsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceVoiceSettingsName** | **[]string** |  | 
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


## DevicevoicesettingsGet

> DevicevoicesettingsGet(ctx).DeviceVoiceSettingsName(deviceVoiceSettingsName).IncludeData(includeData).Execute()

Get all Device Voice Settings



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
	deviceVoiceSettingsName := "deviceVoiceSettingsName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(context.Background()).DeviceVoiceSettingsName(deviceVoiceSettingsName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceVoiceSettingsAPI.DevicevoicesettingsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicevoicesettingsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceVoiceSettingsName** | **string** |  | 
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


## DevicevoicesettingsPatch

> DevicevoicesettingsPatch(ctx).ChangesetName(changesetName).DevicevoicesettingsPutRequest(devicevoicesettingsPutRequest).Execute()

Update Device Voice Setting



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
	devicevoicesettingsPutRequest := *openapiclient.NewDevicevoicesettingsPutRequest() // DevicevoicesettingsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceVoiceSettingsAPI.DevicevoicesettingsPatch(context.Background()).ChangesetName(changesetName).DevicevoicesettingsPutRequest(devicevoicesettingsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceVoiceSettingsAPI.DevicevoicesettingsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicevoicesettingsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **devicevoicesettingsPutRequest** | [**DevicevoicesettingsPutRequest**](DevicevoicesettingsPutRequest.md) |  | 

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


## DevicevoicesettingsPut

> DevicevoicesettingsPut(ctx).ChangesetName(changesetName).DevicevoicesettingsPutRequest(devicevoicesettingsPutRequest).Execute()

Create Device Voice Setting



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
	devicevoicesettingsPutRequest := *openapiclient.NewDevicevoicesettingsPutRequest() // DevicevoicesettingsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceVoiceSettingsAPI.DevicevoicesettingsPut(context.Background()).ChangesetName(changesetName).DevicevoicesettingsPutRequest(devicevoicesettingsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceVoiceSettingsAPI.DevicevoicesettingsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDevicevoicesettingsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **devicevoicesettingsPutRequest** | [**DevicevoicesettingsPutRequest**](DevicevoicesettingsPutRequest.md) |  | 

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

