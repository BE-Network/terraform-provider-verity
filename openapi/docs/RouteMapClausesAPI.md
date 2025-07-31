# \RouteMapClausesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RoutemapclausesDelete**](RouteMapClausesAPI.md#RoutemapclausesDelete) | **Delete** /routemapclauses | Delete Route Map Clause
[**RoutemapclausesGet**](RouteMapClausesAPI.md#RoutemapclausesGet) | **Get** /routemapclauses | Get all Route Map Clauses
[**RoutemapclausesPatch**](RouteMapClausesAPI.md#RoutemapclausesPatch) | **Patch** /routemapclauses | Update Route Map Clause
[**RoutemapclausesPut**](RouteMapClausesAPI.md#RoutemapclausesPut) | **Put** /routemapclauses | Create Route Map Clause



## RoutemapclausesDelete

> RoutemapclausesDelete(ctx).RouteMapClauseName(routeMapClauseName).ChangesetName(changesetName).Execute()

Delete Route Map Clause



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
	routeMapClauseName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RouteMapClausesAPI.RoutemapclausesDelete(context.Background()).RouteMapClauseName(routeMapClauseName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RouteMapClausesAPI.RoutemapclausesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRoutemapclausesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **routeMapClauseName** | **[]string** |  | 
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


## RoutemapclausesGet

> RoutemapclausesGet(ctx).RouteMapClauseName(routeMapClauseName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Route Map Clauses



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
	routeMapClauseName := "routeMapClauseName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RouteMapClausesAPI.RoutemapclausesGet(context.Background()).RouteMapClauseName(routeMapClauseName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RouteMapClausesAPI.RoutemapclausesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRoutemapclausesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **routeMapClauseName** | **string** |  | 
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


## RoutemapclausesPatch

> RoutemapclausesPatch(ctx).ChangesetName(changesetName).RoutemapclausesPutRequest(routemapclausesPutRequest).Execute()

Update Route Map Clause



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
	routemapclausesPutRequest := *openapiclient.NewRoutemapclausesPutRequest() // RoutemapclausesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RouteMapClausesAPI.RoutemapclausesPatch(context.Background()).ChangesetName(changesetName).RoutemapclausesPutRequest(routemapclausesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RouteMapClausesAPI.RoutemapclausesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRoutemapclausesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **routemapclausesPutRequest** | [**RoutemapclausesPutRequest**](RoutemapclausesPutRequest.md) |  | 

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


## RoutemapclausesPut

> RoutemapclausesPut(ctx).ChangesetName(changesetName).RoutemapclausesPutRequest(routemapclausesPutRequest).Execute()

Create Route Map Clause



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
	routemapclausesPutRequest := *openapiclient.NewRoutemapclausesPutRequest() // RoutemapclausesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RouteMapClausesAPI.RoutemapclausesPut(context.Background()).ChangesetName(changesetName).RoutemapclausesPutRequest(routemapclausesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RouteMapClausesAPI.RoutemapclausesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRoutemapclausesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **routemapclausesPutRequest** | [**RoutemapclausesPutRequest**](RoutemapclausesPutRequest.md) |  | 

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

