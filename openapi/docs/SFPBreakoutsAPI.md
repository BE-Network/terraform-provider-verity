# \SFPBreakoutsAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**SfpbreakoutsGet**](SFPBreakoutsAPI.md#SfpbreakoutsGet) | **Get** /sfpbreakouts | Get all SFP Breakouts
[**SfpbreakoutsPatch**](SFPBreakoutsAPI.md#SfpbreakoutsPatch) | **Patch** /sfpbreakouts | Update SFP Breakout



## SfpbreakoutsGet

> SfpbreakoutsGet(ctx).SfpBreakoutsName(sfpBreakoutsName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all SFP Breakouts



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
	sfpBreakoutsName := "sfpBreakoutsName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SFPBreakoutsAPI.SfpbreakoutsGet(context.Background()).SfpBreakoutsName(sfpBreakoutsName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SFPBreakoutsAPI.SfpbreakoutsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSfpbreakoutsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **sfpBreakoutsName** | **string** |  | 
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


## SfpbreakoutsPatch

> SfpbreakoutsPatch(ctx).ChangesetName(changesetName).SfpbreakoutsPatchRequest(sfpbreakoutsPatchRequest).Execute()

Update SFP Breakout



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
	sfpbreakoutsPatchRequest := *openapiclient.NewSfpbreakoutsPatchRequest() // SfpbreakoutsPatchRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.SFPBreakoutsAPI.SfpbreakoutsPatch(context.Background()).ChangesetName(changesetName).SfpbreakoutsPatchRequest(sfpbreakoutsPatchRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SFPBreakoutsAPI.SfpbreakoutsPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSfpbreakoutsPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **sfpbreakoutsPatchRequest** | [**SfpbreakoutsPatchRequest**](SfpbreakoutsPatchRequest.md) |  | 

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

