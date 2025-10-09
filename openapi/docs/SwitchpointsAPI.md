# \SwitchpointsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SwitchpointsCurrentconfigGet**](SwitchpointsAPI.md#SwitchpointsCurrentconfigGet) | **Get** /switchpoints/currentconfig | Get all Switchpoint current configs
[**SwitchpointsDelete**](SwitchpointsAPI.md#SwitchpointsDelete) | **Delete** /switchpoints | Delete Switchpoint
[**SwitchpointsGet**](SwitchpointsAPI.md#SwitchpointsGet) | **Get** /switchpoints | Get all Switchpoints
[**SwitchpointsMarkoutofserviceGet**](SwitchpointsAPI.md#SwitchpointsMarkoutofserviceGet) | **Get** /switchpoints/markoutofservice | Get all marked out of service Switchpoint names
[**SwitchpointsMarkoutofservicePut**](SwitchpointsAPI.md#SwitchpointsMarkoutofservicePut) | **Put** /switchpoints/markoutofservice | Mark switchpoints out of service or back in service
[**SwitchpointsPatch**](SwitchpointsAPI.md#SwitchpointsPatch) | **Patch** /switchpoints | Update Switchpoint
[**SwitchpointsPut**](SwitchpointsAPI.md#SwitchpointsPut) | **Put** /switchpoints | Create Switchpoint
[**SwitchpointsUpgradePatch**](SwitchpointsAPI.md#SwitchpointsUpgradePatch) | **Patch** /switchpoints/upgrade | Update Switchpoint firmware version



## SwitchpointsCurrentconfigGet

> SwitchpointsCurrentconfigGet(ctx).SwitchpointName(switchpointName).Execute()

Get all Switchpoint current configs



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
	switchpointName := "switchpointName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsCurrentconfigGet(context.Background()).SwitchpointName(switchpointName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsCurrentconfigGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsCurrentconfigGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchpointName** | **string** |  | 

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


## SwitchpointsDelete

> SwitchpointsDelete(ctx).SwitchpointName(switchpointName).ChangesetName(changesetName).Execute()

Delete Switchpoint



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
	switchpointName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsDelete(context.Background()).SwitchpointName(switchpointName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchpointName** | **[]string** |  | 
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


## SwitchpointsGet

> SwitchpointsGet(ctx).SwitchpointName(switchpointName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Switchpoints



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
	switchpointName := "switchpointName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsGet(context.Background()).SwitchpointName(switchpointName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchpointName** | **string** |  | 
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


## SwitchpointsMarkoutofserviceGet

> SwitchpointsMarkoutofserviceGet(ctx).Mos(mos).Execute()

Get all marked out of service Switchpoint names



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
	mos := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsMarkoutofserviceGet(context.Background()).Mos(mos).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsMarkoutofserviceGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsMarkoutofserviceGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mos** | **bool** |  | 

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


## SwitchpointsMarkoutofservicePut

> SwitchpointsMarkoutofservicePut(ctx).SwitchpointsMarkoutofservicePutRequest(switchpointsMarkoutofservicePutRequest).Execute()

Mark switchpoints out of service or back in service



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
	switchpointsMarkoutofservicePutRequest := *openapiclient.NewSwitchpointsMarkoutofservicePutRequest([]string{"DeviceNames_example"}) // SwitchpointsMarkoutofservicePutRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsMarkoutofservicePut(context.Background()).SwitchpointsMarkoutofservicePutRequest(switchpointsMarkoutofservicePutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsMarkoutofservicePut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsMarkoutofservicePutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchpointsMarkoutofservicePutRequest** | [**SwitchpointsMarkoutofservicePutRequest**](SwitchpointsMarkoutofservicePutRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[TokenAuth](../README.md#TokenAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SwitchpointsPatch

> SwitchpointsPatch(ctx).ChangesetName(changesetName).SwitchpointsPutRequest(switchpointsPutRequest).Execute()

Update Switchpoint



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
	switchpointsPutRequest := *openapiclient.NewSwitchpointsPutRequest() // SwitchpointsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsPatch(context.Background()).ChangesetName(changesetName).SwitchpointsPutRequest(switchpointsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **switchpointsPutRequest** | [**SwitchpointsPutRequest**](SwitchpointsPutRequest.md) |  | 

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


## SwitchpointsPut

> SwitchpointsPut(ctx).ChangesetName(changesetName).SwitchpointsPutRequest(switchpointsPutRequest).Execute()

Create Switchpoint



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
	switchpointsPutRequest := *openapiclient.NewSwitchpointsPutRequest() // SwitchpointsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsPut(context.Background()).ChangesetName(changesetName).SwitchpointsPutRequest(switchpointsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **switchpointsPutRequest** | [**SwitchpointsPutRequest**](SwitchpointsPutRequest.md) |  | 

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


## SwitchpointsUpgradePatch

> SwitchpointsUpgradePatch(ctx).SwitchpointsUpgradePatchRequest(switchpointsUpgradePatchRequest).Execute()

Update Switchpoint firmware version



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
	switchpointsUpgradePatchRequest := *openapiclient.NewSwitchpointsUpgradePatchRequest("PackageVersion_example", []string{"DeviceNames_example"}) // SwitchpointsUpgradePatchRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SwitchpointsAPI.SwitchpointsUpgradePatch(context.Background()).SwitchpointsUpgradePatchRequest(switchpointsUpgradePatchRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SwitchpointsAPI.SwitchpointsUpgradePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchpointsUpgradePatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **switchpointsUpgradePatchRequest** | [**SwitchpointsUpgradePatchRequest**](SwitchpointsUpgradePatchRequest.md) |  | 

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

