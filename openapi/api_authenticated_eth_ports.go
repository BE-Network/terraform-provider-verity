/*
Verity API

This application demonstrates the usage of Verity API. 

API version: 2.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
)


// AuthenticatedEthPortsAPIService AuthenticatedEthPortsAPI service
type AuthenticatedEthPortsAPIService service

type ApiAuthenticatedethportsDeleteRequest struct {
	ctx context.Context
	ApiService *AuthenticatedEthPortsAPIService
	authenticatedEthPortName *[]string
	changesetName *string
}

func (r ApiAuthenticatedethportsDeleteRequest) AuthenticatedEthPortName(authenticatedEthPortName []string) ApiAuthenticatedethportsDeleteRequest {
	r.authenticatedEthPortName = &authenticatedEthPortName
	return r
}

func (r ApiAuthenticatedethportsDeleteRequest) ChangesetName(changesetName string) ApiAuthenticatedethportsDeleteRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiAuthenticatedethportsDeleteRequest) Execute() (*http.Response, error) {
	return r.ApiService.AuthenticatedethportsDeleteExecute(r)
}

/*
AuthenticatedethportsDelete Delete Authenticated Eth-Port

Deletes an existing Authenticated Eth-Port from the system if changeset_name is empty, from a changeset if its name is provided.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiAuthenticatedethportsDeleteRequest
*/
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsDelete(ctx context.Context) ApiAuthenticatedethportsDeleteRequest {
	return ApiAuthenticatedethportsDeleteRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsDeleteExecute(r ApiAuthenticatedethportsDeleteRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuthenticatedEthPortsAPIService.AuthenticatedethportsDelete")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/authenticatedethports"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.authenticatedEthPortName == nil {
		return nil, reportError("authenticatedEthPortName is required and must be specified")
	}

	{
		t := *r.authenticatedEthPortName
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "authenticated_eth_port_name", s.Index(i).Interface(), "form", "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "authenticated_eth_port_name", t, "form", "multi")
		}
	}
	if r.changesetName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "changeset_name", r.changesetName, "form", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiAuthenticatedethportsGetRequest struct {
	ctx context.Context
	ApiService *AuthenticatedEthPortsAPIService
	authenticatedEthPortName *string
	includeData *bool
	changesetName *string
}

func (r ApiAuthenticatedethportsGetRequest) AuthenticatedEthPortName(authenticatedEthPortName string) ApiAuthenticatedethportsGetRequest {
	r.authenticatedEthPortName = &authenticatedEthPortName
	return r
}

func (r ApiAuthenticatedethportsGetRequest) IncludeData(includeData bool) ApiAuthenticatedethportsGetRequest {
	r.includeData = &includeData
	return r
}

func (r ApiAuthenticatedethportsGetRequest) ChangesetName(changesetName string) ApiAuthenticatedethportsGetRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiAuthenticatedethportsGetRequest) Execute() (*http.Response, error) {
	return r.ApiService.AuthenticatedethportsGetExecute(r)
}

/*
AuthenticatedethportsGet Get all Authenticated Eth-Ports

Retrieves all Authenticated Eth-Ports from the system.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiAuthenticatedethportsGetRequest
*/
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsGet(ctx context.Context) ApiAuthenticatedethportsGetRequest {
	return ApiAuthenticatedethportsGetRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsGetExecute(r ApiAuthenticatedethportsGetRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuthenticatedEthPortsAPIService.AuthenticatedethportsGet")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/authenticatedethports"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.authenticatedEthPortName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "authenticated_eth_port_name", r.authenticatedEthPortName, "form", "")
	}
	if r.includeData != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "include_data", r.includeData, "form", "")
	}
	if r.changesetName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "changeset_name", r.changesetName, "form", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiAuthenticatedethportsPatchRequest struct {
	ctx context.Context
	ApiService *AuthenticatedEthPortsAPIService
	changesetName *string
	authenticatedethportsPutRequest *AuthenticatedethportsPutRequest
}

func (r ApiAuthenticatedethportsPatchRequest) ChangesetName(changesetName string) ApiAuthenticatedethportsPatchRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiAuthenticatedethportsPatchRequest) AuthenticatedethportsPutRequest(authenticatedethportsPutRequest AuthenticatedethportsPutRequest) ApiAuthenticatedethportsPatchRequest {
	r.authenticatedethportsPutRequest = &authenticatedethportsPutRequest
	return r
}

func (r ApiAuthenticatedethportsPatchRequest) Execute() (*http.Response, error) {
	return r.ApiService.AuthenticatedethportsPatchExecute(r)
}

/*
AuthenticatedethportsPatch Update Authenticated Eth-Port

Update Authenticated Eth-Port into the system if changeset_name is empty, into a changeset if its name is provided.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiAuthenticatedethportsPatchRequest
*/
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsPatch(ctx context.Context) ApiAuthenticatedethportsPatchRequest {
	return ApiAuthenticatedethportsPatchRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsPatchExecute(r ApiAuthenticatedethportsPatchRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuthenticatedEthPortsAPIService.AuthenticatedethportsPatch")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/authenticatedethports"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.changesetName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "changeset_name", r.changesetName, "form", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.authenticatedethportsPutRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiAuthenticatedethportsPutRequest struct {
	ctx context.Context
	ApiService *AuthenticatedEthPortsAPIService
	changesetName *string
	authenticatedethportsPutRequest *AuthenticatedethportsPutRequest
}

func (r ApiAuthenticatedethportsPutRequest) ChangesetName(changesetName string) ApiAuthenticatedethportsPutRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiAuthenticatedethportsPutRequest) AuthenticatedethportsPutRequest(authenticatedethportsPutRequest AuthenticatedethportsPutRequest) ApiAuthenticatedethportsPutRequest {
	r.authenticatedethportsPutRequest = &authenticatedethportsPutRequest
	return r
}

func (r ApiAuthenticatedethportsPutRequest) Execute() (*http.Response, error) {
	return r.ApiService.AuthenticatedethportsPutExecute(r)
}

/*
AuthenticatedethportsPut Create Authenticated Eth-Port

Create Authenticated Eth-Port into the system if changeset_name is empty, into a changeset if its name is provided.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiAuthenticatedethportsPutRequest
*/
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsPut(ctx context.Context) ApiAuthenticatedethportsPutRequest {
	return ApiAuthenticatedethportsPutRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *AuthenticatedEthPortsAPIService) AuthenticatedethportsPutExecute(r ApiAuthenticatedethportsPutRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "AuthenticatedEthPortsAPIService.AuthenticatedethportsPut")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/authenticatedethports"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.changesetName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "changeset_name", r.changesetName, "form", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.authenticatedethportsPutRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}
