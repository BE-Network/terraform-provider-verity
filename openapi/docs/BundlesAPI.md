# \BundlesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BundlesGet**](BundlesAPI.md#BundlesGet) | **Get** /bundles | Get all bundles
[**BundlesPatch**](BundlesAPI.md#BundlesPatch) | **Patch** /bundles | Update bundle



## BundlesGet

> BundlesGet(ctx).BundleName(bundleName).IncludeData(includeData).Execute()

Get all bundles



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
	bundleName := "bundleName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.BundlesAPI.BundlesGet(context.Background()).BundleName(bundleName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BundlesAPI.BundlesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiBundlesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bundleName** | **string** |  | 
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


## BundlesPatch

> BundlesPatch(ctx).ChangesetName(changesetName).BundlesPatchRequest(bundlesPatchRequest).Execute()

Update bundle



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
	bundlesPatchRequest := *openapiclient.NewBundlesPatchRequest() // BundlesPatchRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.BundlesAPI.BundlesPatch(context.Background()).ChangesetName(changesetName).BundlesPatchRequest(bundlesPatchRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BundlesAPI.BundlesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiBundlesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **bundlesPatchRequest** | [**BundlesPatchRequest**](BundlesPatchRequest.md) |  | 

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

