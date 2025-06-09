# \PacketBrokerAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PacketbrokerDelete**](PacketBrokerAPI.md#PacketbrokerDelete) | **Delete** /packetbroker | Delete PB Egress Profile
[**PacketbrokerGet**](PacketBrokerAPI.md#PacketbrokerGet) | **Get** /packetbroker | Get all PB Egress Profiles
[**PacketbrokerPatch**](PacketBrokerAPI.md#PacketbrokerPatch) | **Patch** /packetbroker | Update PB Egress Profile
[**PacketbrokerPut**](PacketBrokerAPI.md#PacketbrokerPut) | **Put** /packetbroker | Create PB Egress Profile



## PacketbrokerDelete

> PacketbrokerDelete(ctx).PbEgressProfileName(pbEgressProfileName).ChangesetName(changesetName).Execute()

Delete PB Egress Profile



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
	pbEgressProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketBrokerAPI.PacketbrokerDelete(context.Background()).PbEgressProfileName(pbEgressProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketBrokerAPI.PacketbrokerDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketbrokerDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pbEgressProfileName** | **[]string** |  | 
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


## PacketbrokerGet

> PacketbrokerGet(ctx).PbEgressProfileName(pbEgressProfileName).IncludeData(includeData).Execute()

Get all PB Egress Profiles



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
	pbEgressProfileName := "pbEgressProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketBrokerAPI.PacketbrokerGet(context.Background()).PbEgressProfileName(pbEgressProfileName).IncludeData(includeData).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketBrokerAPI.PacketbrokerGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketbrokerGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pbEgressProfileName** | **string** |  | 
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


## PacketbrokerPatch

> PacketbrokerPatch(ctx).ChangesetName(changesetName).PacketbrokerPutRequest(packetbrokerPutRequest).Execute()

Update PB Egress Profile



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
	packetbrokerPutRequest := *openapiclient.NewPacketbrokerPutRequest() // PacketbrokerPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketBrokerAPI.PacketbrokerPatch(context.Background()).ChangesetName(changesetName).PacketbrokerPutRequest(packetbrokerPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketBrokerAPI.PacketbrokerPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketbrokerPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **packetbrokerPutRequest** | [**PacketbrokerPutRequest**](PacketbrokerPutRequest.md) |  | 

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


## PacketbrokerPut

> PacketbrokerPut(ctx).ChangesetName(changesetName).PacketbrokerPutRequest(packetbrokerPutRequest).Execute()

Create PB Egress Profile



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
	packetbrokerPutRequest := *openapiclient.NewPacketbrokerPutRequest() // PacketbrokerPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketBrokerAPI.PacketbrokerPut(context.Background()).ChangesetName(changesetName).PacketbrokerPutRequest(packetbrokerPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketBrokerAPI.PacketbrokerPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketbrokerPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **packetbrokerPutRequest** | [**PacketbrokerPutRequest**](PacketbrokerPutRequest.md) |  | 

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

