# \RacksAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RacksDelete**](RacksAPI.md#RacksDelete) | **Delete** /racks | Delete rack
[**RacksGet**](RacksAPI.md#RacksGet) | **Get** /racks | Get all racks
[**RacksPatch**](RacksAPI.md#RacksPatch) | **Patch** /racks | Update rack
[**RacksPut**](RacksAPI.md#RacksPut) | **Put** /racks | Create rack



## RacksDelete

> RacksDelete(ctx).RackName(rackName).ChangesetName(changesetName).Execute()

Delete rack



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
	rackName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RacksAPI.RacksDelete(context.Background()).RackName(rackName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RacksAPI.RacksDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRacksDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **rackName** | **[]string** |  | 
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


## RacksGet

> RacksGet(ctx).RackName(rackName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all racks



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
	rackName := "rackName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RacksAPI.RacksGet(context.Background()).RackName(rackName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RacksAPI.RacksGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRacksGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **rackName** | **string** |  | 
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


## RacksPatch

> RacksPatch(ctx).ChangesetName(changesetName).RacksPutRequest(racksPutRequest).Execute()

Update rack



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
	racksPutRequest := *openapiclient.NewRacksPutRequest() // RacksPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RacksAPI.RacksPatch(context.Background()).ChangesetName(changesetName).RacksPutRequest(racksPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RacksAPI.RacksPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRacksPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **racksPutRequest** | [**RacksPutRequest**](RacksPutRequest.md) |  | 

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


## RacksPut

> RacksPut(ctx).ChangesetName(changesetName).RacksPutRequest(racksPutRequest).Execute()

Create rack



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
	racksPutRequest := *openapiclient.NewRacksPutRequest() // RacksPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RacksAPI.RacksPut(context.Background()).ChangesetName(changesetName).RacksPutRequest(racksPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RacksAPI.RacksPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRacksPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **racksPutRequest** | [**RacksPutRequest**](RacksPutRequest.md) |  | 

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

