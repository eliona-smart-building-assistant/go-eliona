/*
Eliona API

API to access Eliona Smart Building Assistant

API version: 2.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// AssetType A type of assets
type AssetType struct {
	// The unique name for this eliona type
	Name string `json:"name"`
	// Is this a customer created type or not
	Custom bool `json:"custom"`
	// The vendor providing assets of this type
	Vendor *string `json:"vendor,omitempty"`
	// The specific model of assets of this type
	Model       *string      `json:"model,omitempty"`
	Translation *Translation `json:"translation,omitempty"`
	// The url describing assets of this type
	Urldoc *string `json:"urldoc,omitempty"`
	// Icon name corresponding to assets of this type
	Icon *string `json:"icon,omitempty"`
	// List of named attributes
	Attributes []Attribute `json:"attributes,omitempty"`
}

// NewAssetType instantiates a new AssetType object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAssetType(name string, custom bool) *AssetType {
	this := AssetType{}
	this.Name = name
	this.Custom = custom
	return &this
}

// NewAssetTypeWithDefaults instantiates a new AssetType object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAssetTypeWithDefaults() *AssetType {
	this := AssetType{}
	var custom bool = true
	this.Custom = custom
	return &this
}

// GetName returns the Name field value
func (o *AssetType) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *AssetType) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *AssetType) SetName(v string) {
	o.Name = v
}

// GetCustom returns the Custom field value
func (o *AssetType) GetCustom() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Custom
}

// GetCustomOk returns a tuple with the Custom field value
// and a boolean to check if the value has been set.
func (o *AssetType) GetCustomOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Custom, true
}

// SetCustom sets field value
func (o *AssetType) SetCustom(v bool) {
	o.Custom = v
}

// GetVendor returns the Vendor field value if set, zero value otherwise.
func (o *AssetType) GetVendor() string {
	if o == nil || o.Vendor == nil {
		var ret string
		return ret
	}
	return *o.Vendor
}

// GetVendorOk returns a tuple with the Vendor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssetType) GetVendorOk() (*string, bool) {
	if o == nil || o.Vendor == nil {
		return nil, false
	}
	return o.Vendor, true
}

// HasVendor returns a boolean if a field has been set.
func (o *AssetType) HasVendor() bool {
	if o != nil && o.Vendor != nil {
		return true
	}

	return false
}

// SetVendor gets a reference to the given string and assigns it to the Vendor field.
func (o *AssetType) SetVendor(v string) {
	o.Vendor = &v
}

// GetModel returns the Model field value if set, zero value otherwise.
func (o *AssetType) GetModel() string {
	if o == nil || o.Model == nil {
		var ret string
		return ret
	}
	return *o.Model
}

// GetModelOk returns a tuple with the Model field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssetType) GetModelOk() (*string, bool) {
	if o == nil || o.Model == nil {
		return nil, false
	}
	return o.Model, true
}

// HasModel returns a boolean if a field has been set.
func (o *AssetType) HasModel() bool {
	if o != nil && o.Model != nil {
		return true
	}

	return false
}

// SetModel gets a reference to the given string and assigns it to the Model field.
func (o *AssetType) SetModel(v string) {
	o.Model = &v
}

// GetTranslation returns the Translation field value if set, zero value otherwise.
func (o *AssetType) GetTranslation() Translation {
	if o == nil || o.Translation == nil {
		var ret Translation
		return ret
	}
	return *o.Translation
}

// GetTranslationOk returns a tuple with the Translation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssetType) GetTranslationOk() (*Translation, bool) {
	if o == nil || o.Translation == nil {
		return nil, false
	}
	return o.Translation, true
}

// HasTranslation returns a boolean if a field has been set.
func (o *AssetType) HasTranslation() bool {
	if o != nil && o.Translation != nil {
		return true
	}

	return false
}

// SetTranslation gets a reference to the given Translation and assigns it to the Translation field.
func (o *AssetType) SetTranslation(v Translation) {
	o.Translation = &v
}

// GetUrldoc returns the Urldoc field value if set, zero value otherwise.
func (o *AssetType) GetUrldoc() string {
	if o == nil || o.Urldoc == nil {
		var ret string
		return ret
	}
	return *o.Urldoc
}

// GetUrldocOk returns a tuple with the Urldoc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssetType) GetUrldocOk() (*string, bool) {
	if o == nil || o.Urldoc == nil {
		return nil, false
	}
	return o.Urldoc, true
}

// HasUrldoc returns a boolean if a field has been set.
func (o *AssetType) HasUrldoc() bool {
	if o != nil && o.Urldoc != nil {
		return true
	}

	return false
}

// SetUrldoc gets a reference to the given string and assigns it to the Urldoc field.
func (o *AssetType) SetUrldoc(v string) {
	o.Urldoc = &v
}

// GetIcon returns the Icon field value if set, zero value otherwise.
func (o *AssetType) GetIcon() string {
	if o == nil || o.Icon == nil {
		var ret string
		return ret
	}
	return *o.Icon
}

// GetIconOk returns a tuple with the Icon field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssetType) GetIconOk() (*string, bool) {
	if o == nil || o.Icon == nil {
		return nil, false
	}
	return o.Icon, true
}

// HasIcon returns a boolean if a field has been set.
func (o *AssetType) HasIcon() bool {
	if o != nil && o.Icon != nil {
		return true
	}

	return false
}

// SetIcon gets a reference to the given string and assigns it to the Icon field.
func (o *AssetType) SetIcon(v string) {
	o.Icon = &v
}

// GetAttributes returns the Attributes field value if set, zero value otherwise.
func (o *AssetType) GetAttributes() []Attribute {
	if o == nil || o.Attributes == nil {
		var ret []Attribute
		return ret
	}
	return o.Attributes
}

// GetAttributesOk returns a tuple with the Attributes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssetType) GetAttributesOk() ([]Attribute, bool) {
	if o == nil || o.Attributes == nil {
		return nil, false
	}
	return o.Attributes, true
}

// HasAttributes returns a boolean if a field has been set.
func (o *AssetType) HasAttributes() bool {
	if o != nil && o.Attributes != nil {
		return true
	}

	return false
}

// SetAttributes gets a reference to the given []Attribute and assigns it to the Attributes field.
func (o *AssetType) SetAttributes(v []Attribute) {
	o.Attributes = v
}

func (o AssetType) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["name"] = o.Name
	}
	if true {
		toSerialize["custom"] = o.Custom
	}
	if o.Vendor != nil {
		toSerialize["vendor"] = o.Vendor
	}
	if o.Model != nil {
		toSerialize["model"] = o.Model
	}
	if o.Translation != nil {
		toSerialize["translation"] = o.Translation
	}
	if o.Urldoc != nil {
		toSerialize["urldoc"] = o.Urldoc
	}
	if o.Icon != nil {
		toSerialize["icon"] = o.Icon
	}
	if o.Attributes != nil {
		toSerialize["attributes"] = o.Attributes
	}
	return json.Marshal(toSerialize)
}

type NullableAssetType struct {
	value *AssetType
	isSet bool
}

func (v NullableAssetType) Get() *AssetType {
	return v.value
}

func (v *NullableAssetType) Set(val *AssetType) {
	v.value = val
	v.isSet = true
}

func (v NullableAssetType) IsSet() bool {
	return v.isSet
}

func (v *NullableAssetType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAssetType(val *AssetType) *NullableAssetType {
	return &NullableAssetType{value: val, isSet: true}
}

func (v NullableAssetType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAssetType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
