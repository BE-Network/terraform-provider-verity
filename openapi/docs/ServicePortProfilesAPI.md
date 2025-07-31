# \ServicePortProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ServiceportprofilesDelete**](ServicePortProfilesAPI.md#ServiceportprofilesDelete) | **Delete** /serviceportprofiles | Delete Service Port Profile
[**ServiceportprofilesGet**](ServicePortProfilesAPI.md#ServiceportprofilesGet) | **Get** /serviceportprofiles | Get all Service Port Profiles
[**ServiceportprofilesPatch**](ServicePortProfilesAPI.md#ServiceportprofilesPatch) | **Patch** /serviceportprofiles | Update Service Port Profile
[**ServiceportprofilesPut**](ServicePortProfilesAPI.md#ServiceportprofilesPut) | **Put** /serviceportprofiles | Create Service Port Profile



## ServiceportprofilesDelete

> ServiceportprofilesDelete(ctx).ServicePortProfileName(servicePortProfileName).ChangesetName(changesetName).Execute()

Delete Service Port Profile



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
	servicePortProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ServicePortProfilesAPI.ServiceportprofilesDelete(context.Background()).ServicePortProfileName(servicePortProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePortProfilesAPI.ServiceportprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServiceportprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **servicePortProfileName** | **[]string** |  | 
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


## ServiceportprofilesGet

> ServiceportprofilesGet(ctx).ServicePortProfileName(servicePortProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Service Port Profiles



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
	servicePortProfileName := "servicePortProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ServicePortProfilesAPI.ServiceportprofilesGet(context.Background()).ServicePortProfileName(servicePortProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePortProfilesAPI.ServiceportprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServiceportprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **servicePortProfileName** | **string** |  | 
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


## ServiceportprofilesPatch

> ServiceportprofilesPatch(ctx).ChangesetName(changesetName).ServiceportprofilesPutRequest(serviceportprofilesPutRequest).Execute()

Update Service Port Profile



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
	serviceportprofilesPutRequest := *openapiclient.NewServiceportprofilesPutRequest() // ServiceportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ServicePortProfilesAPI.ServiceportprofilesPatch(context.Background()).ChangesetName(changesetName).ServiceportprofilesPutRequest(serviceportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePortProfilesAPI.ServiceportprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServiceportprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **serviceportprofilesPutRequest** | [**ServiceportprofilesPutRequest**](ServiceportprofilesPutRequest.md) |  | 

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


## ServiceportprofilesPut

> ServiceportprofilesPut(ctx).ChangesetName(changesetName).ServiceportprofilesPutRequest(serviceportprofilesPutRequest).Execute()

Create Service Port Profile



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
	serviceportprofilesPutRequest := *openapiclient.NewServiceportprofilesPutRequest() // ServiceportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ServicePortProfilesAPI.ServiceportprofilesPut(context.Background()).ChangesetName(changesetName).ServiceportprofilesPutRequest(serviceportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServicePortProfilesAPI.ServiceportprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiServiceportprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **serviceportprofilesPutRequest** | [**ServiceportprofilesPutRequest**](ServiceportprofilesPutRequest.md) |  | 

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

