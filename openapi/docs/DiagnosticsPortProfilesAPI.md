# \DiagnosticsPortProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DiagnosticsportprofilesDelete**](DiagnosticsPortProfilesAPI.md#DiagnosticsportprofilesDelete) | **Delete** /diagnosticsportprofiles | Delete Diagnostics Port Profile
[**DiagnosticsportprofilesGet**](DiagnosticsPortProfilesAPI.md#DiagnosticsportprofilesGet) | **Get** /diagnosticsportprofiles | Get all Diagnostics Port Profiles
[**DiagnosticsportprofilesPatch**](DiagnosticsPortProfilesAPI.md#DiagnosticsportprofilesPatch) | **Patch** /diagnosticsportprofiles | Update Diagnostics Port Profile
[**DiagnosticsportprofilesPut**](DiagnosticsPortProfilesAPI.md#DiagnosticsportprofilesPut) | **Put** /diagnosticsportprofiles | Create Diagnostics Port Profile



## DiagnosticsportprofilesDelete

> DiagnosticsportprofilesDelete(ctx).DiagnosticsPortProfileName(diagnosticsPortProfileName).ChangesetName(changesetName).Execute()

Delete Diagnostics Port Profile



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
	diagnosticsPortProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesDelete(context.Background()).DiagnosticsPortProfileName(diagnosticsPortProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsPortProfilesAPI.DiagnosticsportprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsportprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **diagnosticsPortProfileName** | **[]string** |  | 
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


## DiagnosticsportprofilesGet

> DiagnosticsportprofilesGet(ctx).DiagnosticsPortProfileName(diagnosticsPortProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Diagnostics Port Profiles



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
	diagnosticsPortProfileName := "diagnosticsPortProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(context.Background()).DiagnosticsPortProfileName(diagnosticsPortProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsportprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **diagnosticsPortProfileName** | **string** |  | 
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


## DiagnosticsportprofilesPatch

> DiagnosticsportprofilesPatch(ctx).ChangesetName(changesetName).DiagnosticsportprofilesPutRequest(diagnosticsportprofilesPutRequest).Execute()

Update Diagnostics Port Profile



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
	diagnosticsportprofilesPutRequest := *openapiclient.NewDiagnosticsportprofilesPutRequest() // DiagnosticsportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPatch(context.Background()).ChangesetName(changesetName).DiagnosticsportprofilesPutRequest(diagnosticsportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsportprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **diagnosticsportprofilesPutRequest** | [**DiagnosticsportprofilesPutRequest**](DiagnosticsportprofilesPutRequest.md) |  | 

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


## DiagnosticsportprofilesPut

> DiagnosticsportprofilesPut(ctx).ChangesetName(changesetName).DiagnosticsportprofilesPutRequest(diagnosticsportprofilesPutRequest).Execute()

Create Diagnostics Port Profile



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
	diagnosticsportprofilesPutRequest := *openapiclient.NewDiagnosticsportprofilesPutRequest() // DiagnosticsportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPut(context.Background()).ChangesetName(changesetName).DiagnosticsportprofilesPutRequest(diagnosticsportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsportprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **diagnosticsportprofilesPutRequest** | [**DiagnosticsportprofilesPutRequest**](DiagnosticsportprofilesPutRequest.md) |  | 

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

