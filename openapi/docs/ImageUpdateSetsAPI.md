# \ImageUpdateSetsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ImageupdatesetsDelete**](ImageUpdateSetsAPI.md#ImageupdatesetsDelete) | **Delete** /imageupdatesets | Delete Image Update Set
[**ImageupdatesetsGet**](ImageUpdateSetsAPI.md#ImageupdatesetsGet) | **Get** /imageupdatesets | Get all Image Update Sets
[**ImageupdatesetsPatch**](ImageUpdateSetsAPI.md#ImageupdatesetsPatch) | **Patch** /imageupdatesets | Update Image Update Set
[**ImageupdatesetsPut**](ImageUpdateSetsAPI.md#ImageupdatesetsPut) | **Put** /imageupdatesets | Create Image Update Set



## ImageupdatesetsDelete

> ImageupdatesetsDelete(ctx).ImageUpdateSetName(imageUpdateSetName).ChangesetName(changesetName).Execute()

Delete Image Update Set



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
	imageUpdateSetName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ImageUpdateSetsAPI.ImageupdatesetsDelete(context.Background()).ImageUpdateSetName(imageUpdateSetName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageUpdateSetsAPI.ImageupdatesetsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiImageupdatesetsDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **imageUpdateSetName** | **[]string** |  | 
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


## ImageupdatesetsGet

> ImageupdatesetsGet(ctx).ImageUpdateSetName(imageUpdateSetName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Image Update Sets



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
	imageUpdateSetName := "imageUpdateSetName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ImageUpdateSetsAPI.ImageupdatesetsGet(context.Background()).ImageUpdateSetName(imageUpdateSetName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageUpdateSetsAPI.ImageupdatesetsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiImageupdatesetsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **imageUpdateSetName** | **string** |  | 
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


## ImageupdatesetsPatch

> ImageupdatesetsPatch(ctx).ChangesetName(changesetName).ImageupdatesetsPutRequest(imageupdatesetsPutRequest).Execute()

Update Image Update Set



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
	imageupdatesetsPutRequest := *openapiclient.NewImageupdatesetsPutRequest() // ImageupdatesetsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ImageUpdateSetsAPI.ImageupdatesetsPatch(context.Background()).ChangesetName(changesetName).ImageupdatesetsPutRequest(imageupdatesetsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageUpdateSetsAPI.ImageupdatesetsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiImageupdatesetsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **imageupdatesetsPutRequest** | [**ImageupdatesetsPutRequest**](ImageupdatesetsPutRequest.md) |  | 

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


## ImageupdatesetsPut

> ImageupdatesetsPut(ctx).ChangesetName(changesetName).ImageupdatesetsPutRequest(imageupdatesetsPutRequest).Execute()

Create Image Update Set



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
	imageupdatesetsPutRequest := *openapiclient.NewImageupdatesetsPutRequest() // ImageupdatesetsPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ImageUpdateSetsAPI.ImageupdatesetsPut(context.Background()).ChangesetName(changesetName).ImageupdatesetsPutRequest(imageupdatesetsPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ImageUpdateSetsAPI.ImageupdatesetsPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiImageupdatesetsPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **imageupdatesetsPutRequest** | [**ImageupdatesetsPutRequest**](ImageupdatesetsPutRequest.md) |  | 

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

