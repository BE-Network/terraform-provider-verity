# \VoicePortProfilesAPI

All URIs are relative to *http://localhost/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**VoiceportprofilesDelete**](VoicePortProfilesAPI.md#VoiceportprofilesDelete) | **Delete** /voiceportprofiles | Delete Voice-Port Profile
[**VoiceportprofilesGet**](VoicePortProfilesAPI.md#VoiceportprofilesGet) | **Get** /voiceportprofiles | Get all Voice-Port Profiles
[**VoiceportprofilesPatch**](VoicePortProfilesAPI.md#VoiceportprofilesPatch) | **Patch** /voiceportprofiles | Update Voice-Port Profile
[**VoiceportprofilesPut**](VoicePortProfilesAPI.md#VoiceportprofilesPut) | **Put** /voiceportprofiles | Create Voice-Port ProfileVoice-Port Profiles



## VoiceportprofilesDelete

> VoiceportprofilesDelete(ctx).VoicePortProfileName(voicePortProfileName).ChangesetName(changesetName).Execute()

Delete Voice-Port Profile



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
	voicePortProfileName := []string{"Inner_example"} // []string | 
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VoicePortProfilesAPI.VoiceportprofilesDelete(context.Background()).VoicePortProfileName(voicePortProfileName).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VoicePortProfilesAPI.VoiceportprofilesDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVoiceportprofilesDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **voicePortProfileName** | **[]string** |  | 
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


## VoiceportprofilesGet

> VoiceportprofilesGet(ctx).VoicePortProfileName(voicePortProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()

Get all Voice-Port Profiles



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
	voicePortProfileName := "voicePortProfileName_example" // string |  (optional)
	includeData := true // bool |  (optional)
	changesetName := "changesetName_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VoicePortProfilesAPI.VoiceportprofilesGet(context.Background()).VoicePortProfileName(voicePortProfileName).IncludeData(includeData).ChangesetName(changesetName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VoicePortProfilesAPI.VoiceportprofilesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVoiceportprofilesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **voicePortProfileName** | **string** |  | 
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


## VoiceportprofilesPatch

> VoiceportprofilesPatch(ctx).ChangesetName(changesetName).VoiceportprofilesPutRequest(voiceportprofilesPutRequest).Execute()

Update Voice-Port Profile



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
	voiceportprofilesPutRequest := *openapiclient.NewVoiceportprofilesPutRequest() // VoiceportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VoicePortProfilesAPI.VoiceportprofilesPatch(context.Background()).ChangesetName(changesetName).VoiceportprofilesPutRequest(voiceportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VoicePortProfilesAPI.VoiceportprofilesPatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVoiceportprofilesPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **voiceportprofilesPutRequest** | [**VoiceportprofilesPutRequest**](VoiceportprofilesPutRequest.md) |  | 

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


## VoiceportprofilesPut

> VoiceportprofilesPut(ctx).ChangesetName(changesetName).VoiceportprofilesPutRequest(voiceportprofilesPutRequest).Execute()

Create Voice-Port ProfileVoice-Port Profiles



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
	voiceportprofilesPutRequest := *openapiclient.NewVoiceportprofilesPutRequest() // VoiceportprofilesPutRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.VoicePortProfilesAPI.VoiceportprofilesPut(context.Background()).ChangesetName(changesetName).VoiceportprofilesPutRequest(voiceportprofilesPutRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `VoicePortProfilesAPI.VoiceportprofilesPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiVoiceportprofilesPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **changesetName** | **string** |  | 
 **voiceportprofilesPutRequest** | [**VoiceportprofilesPutRequest**](VoiceportprofilesPutRequest.md) |  | 

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

