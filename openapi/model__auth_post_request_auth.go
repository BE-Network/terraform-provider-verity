/*
Verity API

This application demonstrates the usage of Verity API. 

API version: 2.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the AuthPostRequestAuth type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthPostRequestAuth{}

// AuthPostRequestAuth struct for AuthPostRequestAuth
type AuthPostRequestAuth struct {
	// The username for authentication
	Username string `json:"username"`
	// The password for authentication
	Password string `json:"password"`
}

type _AuthPostRequestAuth AuthPostRequestAuth

// NewAuthPostRequestAuth instantiates a new AuthPostRequestAuth object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthPostRequestAuth(username string, password string) *AuthPostRequestAuth {
	this := AuthPostRequestAuth{}
	this.Username = username
	this.Password = password
	return &this
}

// NewAuthPostRequestAuthWithDefaults instantiates a new AuthPostRequestAuth object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthPostRequestAuthWithDefaults() *AuthPostRequestAuth {
	this := AuthPostRequestAuth{}
	return &this
}

// GetUsername returns the Username field value
func (o *AuthPostRequestAuth) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *AuthPostRequestAuth) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *AuthPostRequestAuth) SetUsername(v string) {
	o.Username = v
}

// GetPassword returns the Password field value
func (o *AuthPostRequestAuth) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *AuthPostRequestAuth) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *AuthPostRequestAuth) SetPassword(v string) {
	o.Password = v
}

func (o AuthPostRequestAuth) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthPostRequestAuth) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["username"] = o.Username
	toSerialize["password"] = o.Password
	return toSerialize, nil
}

func (o *AuthPostRequestAuth) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"username",
		"password",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varAuthPostRequestAuth := _AuthPostRequestAuth{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varAuthPostRequestAuth)

	if err != nil {
		return err
	}

	*o = AuthPostRequestAuth(varAuthPostRequestAuth)

	return err
}

type NullableAuthPostRequestAuth struct {
	value *AuthPostRequestAuth
	isSet bool
}

func (v NullableAuthPostRequestAuth) Get() *AuthPostRequestAuth {
	return v.value
}

func (v *NullableAuthPostRequestAuth) Set(val *AuthPostRequestAuth) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthPostRequestAuth) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthPostRequestAuth) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthPostRequestAuth(val *AuthPostRequestAuth) *NullableAuthPostRequestAuth {
	return &NullableAuthPostRequestAuth{value: val, isSet: true}
}

func (v NullableAuthPostRequestAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthPostRequestAuth) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


