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


// DeviceControllersAPIService DeviceControllersAPI service
type DeviceControllersAPIService service

type ApiDevicecontrollersDeleteRequest struct {
	ctx context.Context
	ApiService *DeviceControllersAPIService
	deviceControllerName *[]string
	changesetName *string
}

func (r ApiDevicecontrollersDeleteRequest) DeviceControllerName(deviceControllerName []string) ApiDevicecontrollersDeleteRequest {
	r.deviceControllerName = &deviceControllerName
	return r
}

func (r ApiDevicecontrollersDeleteRequest) ChangesetName(changesetName string) ApiDevicecontrollersDeleteRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiDevicecontrollersDeleteRequest) Execute() (*http.Response, error) {
	return r.ApiService.DevicecontrollersDeleteExecute(r)
}

/*
DevicecontrollersDelete Delete Device Controllers

Deletes an existing Device Controllers from the system if changeset_name is empty, from a changeset if its name is provided.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiDevicecontrollersDeleteRequest
*/
func (a *DeviceControllersAPIService) DevicecontrollersDelete(ctx context.Context) ApiDevicecontrollersDeleteRequest {
	return ApiDevicecontrollersDeleteRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *DeviceControllersAPIService) DevicecontrollersDeleteExecute(r ApiDevicecontrollersDeleteRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DeviceControllersAPIService.DevicecontrollersDelete")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/devicecontrollers"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.deviceControllerName == nil {
		return nil, reportError("deviceControllerName is required and must be specified")
	}

	{
		t := *r.deviceControllerName
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				parameterAddToHeaderOrQuery(localVarQueryParams, "device_controller_name", s.Index(i).Interface(), "form", "multi")
			}
		} else {
			parameterAddToHeaderOrQuery(localVarQueryParams, "device_controller_name", t, "form", "multi")
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

type ApiDevicecontrollersGetRequest struct {
	ctx context.Context
	ApiService *DeviceControllersAPIService
	deviceControllerName *string
	includeData *bool
	changesetName *string
}

func (r ApiDevicecontrollersGetRequest) DeviceControllerName(deviceControllerName string) ApiDevicecontrollersGetRequest {
	r.deviceControllerName = &deviceControllerName
	return r
}

func (r ApiDevicecontrollersGetRequest) IncludeData(includeData bool) ApiDevicecontrollersGetRequest {
	r.includeData = &includeData
	return r
}

func (r ApiDevicecontrollersGetRequest) ChangesetName(changesetName string) ApiDevicecontrollersGetRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiDevicecontrollersGetRequest) Execute() (*http.Response, error) {
	return r.ApiService.DevicecontrollersGetExecute(r)
}

/*
DevicecontrollersGet Get all Device Controllers

Retrieves all Device Controllers from the system.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiDevicecontrollersGetRequest
*/
func (a *DeviceControllersAPIService) DevicecontrollersGet(ctx context.Context) ApiDevicecontrollersGetRequest {
	return ApiDevicecontrollersGetRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *DeviceControllersAPIService) DevicecontrollersGetExecute(r ApiDevicecontrollersGetRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DeviceControllersAPIService.DevicecontrollersGet")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/devicecontrollers"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.deviceControllerName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "device_controller_name", r.deviceControllerName, "form", "")
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

type ApiDevicecontrollersPatchRequest struct {
	ctx context.Context
	ApiService *DeviceControllersAPIService
	changesetName *string
	devicecontrollersPutRequest *DevicecontrollersPutRequest
}

func (r ApiDevicecontrollersPatchRequest) ChangesetName(changesetName string) ApiDevicecontrollersPatchRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiDevicecontrollersPatchRequest) DevicecontrollersPutRequest(devicecontrollersPutRequest DevicecontrollersPutRequest) ApiDevicecontrollersPatchRequest {
	r.devicecontrollersPutRequest = &devicecontrollersPutRequest
	return r
}

func (r ApiDevicecontrollersPatchRequest) Execute() (*http.Response, error) {
	return r.ApiService.DevicecontrollersPatchExecute(r)
}

/*
DevicecontrollersPatch Update Device Controller

Update Device Controller into the system if changeset_name is empty, into a changeset if its name is provided.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiDevicecontrollersPatchRequest
*/
func (a *DeviceControllersAPIService) DevicecontrollersPatch(ctx context.Context) ApiDevicecontrollersPatchRequest {
	return ApiDevicecontrollersPatchRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *DeviceControllersAPIService) DevicecontrollersPatchExecute(r ApiDevicecontrollersPatchRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPatch
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DeviceControllersAPIService.DevicecontrollersPatch")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/devicecontrollers"

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
	localVarPostBody = r.devicecontrollersPutRequest
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

type ApiDevicecontrollersPutRequest struct {
	ctx context.Context
	ApiService *DeviceControllersAPIService
	changesetName *string
	devicecontrollersPutRequest *DevicecontrollersPutRequest
}

func (r ApiDevicecontrollersPutRequest) ChangesetName(changesetName string) ApiDevicecontrollersPutRequest {
	r.changesetName = &changesetName
	return r
}

func (r ApiDevicecontrollersPutRequest) DevicecontrollersPutRequest(devicecontrollersPutRequest DevicecontrollersPutRequest) ApiDevicecontrollersPutRequest {
	r.devicecontrollersPutRequest = &devicecontrollersPutRequest
	return r
}

func (r ApiDevicecontrollersPutRequest) Execute() (*http.Response, error) {
	return r.ApiService.DevicecontrollersPutExecute(r)
}

/*
DevicecontrollersPut Create Device Controller

Create Device Controller into the system if changeset_name is empty, into a changeset if its name is provided.


 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiDevicecontrollersPutRequest
*/
func (a *DeviceControllersAPIService) DevicecontrollersPut(ctx context.Context) ApiDevicecontrollersPutRequest {
	return ApiDevicecontrollersPutRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
func (a *DeviceControllersAPIService) DevicecontrollersPutExecute(r ApiDevicecontrollersPutRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DeviceControllersAPIService.DevicecontrollersPut")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/devicecontrollers"

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
	localVarPostBody = r.devicecontrollersPutRequest
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
