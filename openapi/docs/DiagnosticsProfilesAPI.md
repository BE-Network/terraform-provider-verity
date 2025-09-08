# \DiagnosticsProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DiagnosticsprofilesDelete**](DiagnosticsProfilesAPI.md#DiagnosticsprofilesDelete) | **Delete** /diagnosticsprofiles | Delete Diagnostics Profile
[**DiagnosticsprofilesGet**](DiagnosticsProfilesAPI.md#DiagnosticsprofilesGet) | **Get** /diagnosticsprofiles | Get all Diagnostics Profiles
[**DiagnosticsprofilesPatch**](DiagnosticsProfilesAPI.md#DiagnosticsprofilesPatch) | **Patch** /diagnosticsprofiles | Update Diagnostics Profile
[**DiagnosticsprofilesPut**](DiagnosticsProfilesAPI.md#DiagnosticsprofilesPut) | **Put** /diagnosticsprofiles | Create Diagnostics Profile



## DiagnosticsprofilesDelete

> DiagnosticsprofilesDelete(ctx).DiagnosticsProfileName(diagnosticsProfileName).ChangesetName(changesetName).Execute()

Delete Diagnostics Profile



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
	diagnosticsProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsProfilesAPI.DiagnosticsprofilesDelete(context.Background()).DiagnosticsProfileName(diagnosticsProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsProfilesAPI.DiagnosticsprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **diagnosticsProfileName** | **[]string** |  | 
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


## DiagnosticsprofilesGet

> DiagnosticsprofilesGet(ctx).DiagnosticsProfileName(diagnosticsProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Diagnostics Profiles



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
	diagnosticsProfileName := "diagnosticsProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(context.Background()).DiagnosticsProfileName(diagnosticsProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsProfilesAPI.DiagnosticsprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **diagnosticsProfileName** | **string** |  | 
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


## DiagnosticsprofilesPatch

> DiagnosticsprofilesPatch(ctx).ChangesetName(changesetName).DiagnosticsprofilesPutRequest(diagnosticsprofilesPutRequest).Execute()

Update Diagnostics Profile



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
	diagnosticsprofilesPutRequest := *openapiclient.NewDiagnosticsprofilesPutRequest() // DiagnosticsprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsProfilesAPI.DiagnosticsprofilesPatch(context.Background()).ChangesetName(changesetName).DiagnosticsprofilesPutRequest(diagnosticsprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsProfilesAPI.DiagnosticsprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **diagnosticsprofilesPutRequest** | [**DiagnosticsprofilesPutRequest**](DiagnosticsprofilesPutRequest.md) |  | 

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


## DiagnosticsprofilesPut

> DiagnosticsprofilesPut(ctx).ChangesetName(changesetName).DiagnosticsprofilesPutRequest(diagnosticsprofilesPutRequest).Execute()

Create Diagnostics Profile



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
	diagnosticsprofilesPutRequest := *openapiclient.NewDiagnosticsprofilesPutRequest() // DiagnosticsprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DiagnosticsProfilesAPI.DiagnosticsprofilesPut(context.Background()).ChangesetName(changesetName).DiagnosticsprofilesPutRequest(diagnosticsprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DiagnosticsProfilesAPI.DiagnosticsprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiagnosticsprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **diagnosticsprofilesPutRequest** | [**DiagnosticsprofilesPutRequest**](DiagnosticsprofilesPutRequest.md) |  | 

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

