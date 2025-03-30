# \EthPortSettingsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EthportsettingsDelete**](EthPortSettingsAPI.md#EthportsettingsDelete) | **Delete** /ethportsettings | Delete Eth-Port Settings
[**EthportsettingsGet**](EthPortSettingsAPI.md#EthportsettingsGet) | **Get** /ethportsettings | Get all Eth-Port Settings
[**EthportsettingsPatch**](EthPortSettingsAPI.md#EthportsettingsPatch) | **Patch** /ethportsettings | Update Eth-Port Settings
[**EthportsettingsPut**](EthPortSettingsAPI.md#EthportsettingsPut) | **Put** /ethportsettings | Create Eth-Port Settings



## EthportsettingsDelete

> EthportsettingsDelete(ctx).PortName(portName).ChangesetName(changesetName).Execute()

Delete Eth-Port Settings



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
	portName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortSettingsAPI.EthportsettingsDelete(context.Background()).PortName(portName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortSettingsAPI.EthportsettingsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportsettingsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **portName** | **[]string** |  | 
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


## EthportsettingsGet

> EthportsettingsGet(ctx).PortName(portName).IncludeData(includeData).Execute()

Get all Eth-Port Settings



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
	portName := "portName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortSettingsAPI.EthportsettingsGet(context.Background()).PortName(portName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortSettingsAPI.EthportsettingsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportsettingsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **portName** | **string** |  | 
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


## EthportsettingsPatch

> EthportsettingsPatch(ctx).ChangesetName(changesetName).EthportsettingsPutRequest(ethportsettingsPutRequest).Execute()

Update Eth-Port Settings



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
	ethportsettingsPutRequest := *openapiclient.NewEthportsettingsPutRequest() // EthportsettingsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortSettingsAPI.EthportsettingsPatch(context.Background()).ChangesetName(changesetName).EthportsettingsPutRequest(ethportsettingsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortSettingsAPI.EthportsettingsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportsettingsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ethportsettingsPutRequest** | [**EthportsettingsPutRequest**](EthportsettingsPutRequest.md) |  | 

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


## EthportsettingsPut

> EthportsettingsPut(ctx).ChangesetName(changesetName).EthportsettingsPutRequest(ethportsettingsPutRequest).Execute()

Create Eth-Port Settings



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
	ethportsettingsPutRequest := *openapiclient.NewEthportsettingsPutRequest() // EthportsettingsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.EthPortSettingsAPI.EthportsettingsPut(context.Background()).ChangesetName(changesetName).EthportsettingsPutRequest(ethportsettingsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EthPortSettingsAPI.EthportsettingsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEthportsettingsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **ethportsettingsPutRequest** | [**EthportsettingsPutRequest**](EthportsettingsPutRequest.md) |  | 

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

