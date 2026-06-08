# \DeviceAAAProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeviceaaaprofilesDelete**](DeviceAAAProfilesAPI.md#DeviceaaaprofilesDelete) | **Delete** /deviceaaaprofiles | Delete Device AAA Profile
[**DeviceaaaprofilesGet**](DeviceAAAProfilesAPI.md#DeviceaaaprofilesGet) | **Get** /deviceaaaprofiles | Get all Device AAA Profiles
[**DeviceaaaprofilesPatch**](DeviceAAAProfilesAPI.md#DeviceaaaprofilesPatch) | **Patch** /deviceaaaprofiles | Update Device AAA Profile
[**DeviceaaaprofilesPut**](DeviceAAAProfilesAPI.md#DeviceaaaprofilesPut) | **Put** /deviceaaaprofiles | Create Device AAA Profile



## DeviceaaaprofilesDelete

> DeviceaaaprofilesDelete(ctx).DeviceAaaProfileName(deviceAaaProfileName).ChangesetName(changesetName).Execute()

Delete Device AAA Profile



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
	deviceAaaProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceAAAProfilesAPI.DeviceaaaprofilesDelete(context.Background()).DeviceAaaProfileName(deviceAaaProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceAAAProfilesAPI.DeviceaaaprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeviceaaaprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceAaaProfileName** | **[]string** |  | 
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


## DeviceaaaprofilesGet

> DeviceaaaprofilesGet(ctx).DeviceAaaProfileName(deviceAaaProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Device AAA Profiles



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
	deviceAaaProfileName := "deviceAaaProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceAAAProfilesAPI.DeviceaaaprofilesGet(context.Background()).DeviceAaaProfileName(deviceAaaProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceAAAProfilesAPI.DeviceaaaprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeviceaaaprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceAaaProfileName** | **string** |  | 
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


## DeviceaaaprofilesPatch

> DeviceaaaprofilesPatch(ctx).ChangesetName(changesetName).DeviceaaaprofilesPutRequest(deviceaaaprofilesPutRequest).Execute()

Update Device AAA Profile



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
	deviceaaaprofilesPutRequest := *openapiclient.NewDeviceaaaprofilesPutRequest() // DeviceaaaprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceAAAProfilesAPI.DeviceaaaprofilesPatch(context.Background()).ChangesetName(changesetName).DeviceaaaprofilesPutRequest(deviceaaaprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceAAAProfilesAPI.DeviceaaaprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeviceaaaprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **deviceaaaprofilesPutRequest** | [**DeviceaaaprofilesPutRequest**](DeviceaaaprofilesPutRequest.md) |  | 

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


## DeviceaaaprofilesPut

> DeviceaaaprofilesPut(ctx).ChangesetName(changesetName).DeviceaaaprofilesPutRequest(deviceaaaprofilesPutRequest).Execute()

Create Device AAA Profile



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
	deviceaaaprofilesPutRequest := *openapiclient.NewDeviceaaaprofilesPutRequest() // DeviceaaaprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DeviceAAAProfilesAPI.DeviceaaaprofilesPut(context.Background()).ChangesetName(changesetName).DeviceaaaprofilesPutRequest(deviceaaaprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DeviceAAAProfilesAPI.DeviceaaaprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeviceaaaprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **deviceaaaprofilesPutRequest** | [**DeviceaaaprofilesPutRequest**](DeviceaaaprofilesPutRequest.md) |  | 

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

