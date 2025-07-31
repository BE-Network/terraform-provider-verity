# \PacketQueuesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PacketqueuesDelete**](PacketQueuesAPI.md#PacketqueuesDelete) | **Delete** /packetqueues | Delete Packet Queue
[**PacketqueuesGet**](PacketQueuesAPI.md#PacketqueuesGet) | **Get** /packetqueues | Get all Packet Queues
[**PacketqueuesPatch**](PacketQueuesAPI.md#PacketqueuesPatch) | **Patch** /packetqueues | Update Packet Queue
[**PacketqueuesPut**](PacketQueuesAPI.md#PacketqueuesPut) | **Put** /packetqueues | Create Packet Queue



## PacketqueuesDelete

> PacketqueuesDelete(ctx).PacketQueueName(packetQueueName).ChangesetName(changesetName).Execute()

Delete Packet Queue



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
	packetQueueName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketQueuesAPI.PacketqueuesDelete(context.Background()).PacketQueueName(packetQueueName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketQueuesAPI.PacketqueuesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketqueuesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **packetQueueName** | **[]string** |  | 
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


## PacketqueuesGet

> PacketqueuesGet(ctx).PacketQueueName(packetQueueName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Packet Queues



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
	packetQueueName := "packetQueueName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketQueuesAPI.PacketqueuesGet(context.Background()).PacketQueueName(packetQueueName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketQueuesAPI.PacketqueuesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketqueuesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **packetQueueName** | **string** |  | 
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


## PacketqueuesPatch

> PacketqueuesPatch(ctx).ChangesetName(changesetName).PacketqueuesPutRequest(packetqueuesPutRequest).Execute()

Update Packet Queue



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
	packetqueuesPutRequest := *openapiclient.NewPacketqueuesPutRequest() // PacketqueuesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketQueuesAPI.PacketqueuesPatch(context.Background()).ChangesetName(changesetName).PacketqueuesPutRequest(packetqueuesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketQueuesAPI.PacketqueuesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketqueuesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **packetqueuesPutRequest** | [**PacketqueuesPutRequest**](PacketqueuesPutRequest.md) |  | 

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


## PacketqueuesPut

> PacketqueuesPut(ctx).ChangesetName(changesetName).PacketqueuesPutRequest(packetqueuesPutRequest).Execute()

Create Packet Queue



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
	packetqueuesPutRequest := *openapiclient.NewPacketqueuesPutRequest() // PacketqueuesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.PacketQueuesAPI.PacketqueuesPut(context.Background()).ChangesetName(changesetName).PacketqueuesPutRequest(packetqueuesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PacketQueuesAPI.PacketqueuesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPacketqueuesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **packetqueuesPutRequest** | [**PacketqueuesPutRequest**](PacketqueuesPutRequest.md) |  | 

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

