# \ImageUpdateSetsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ImageupdatesetsGet**](ImageUpdateSetsAPI.md#ImageupdatesetsGet) | **Get** /imageupdatesets | Get all Image Update Sets
[**ImageupdatesetsPatch**](ImageUpdateSetsAPI.md#ImageupdatesetsPatch) | **Patch** /imageupdatesets | Update Image Update Set



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

> ImageupdatesetsPatch(ctx).ChangesetName(changesetName).ImageupdatesetsPatchRequest(imageupdatesetsPatchRequest).Execute()

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
	imageupdatesetsPatchRequest := *openapiclient.NewImageupdatesetsPatchRequest() // ImageupdatesetsPatchRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.ImageUpdateSetsAPI.ImageupdatesetsPatch(context.Background()).ChangesetName(changesetName).ImageupdatesetsPatchRequest(imageupdatesetsPatchRequest).Execute()
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
 **imageupdatesetsPatchRequest** | [**ImageupdatesetsPatchRequest**](ImageupdatesetsPatchRequest.md) |  | 

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

