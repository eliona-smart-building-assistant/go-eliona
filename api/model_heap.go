/*
Eliona API

API to access Eliona Smart Building Assistant

API version: 2.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
	"time"
)

// Heap Heap data
type Heap struct {
	// ID of the corresponding eliona
	AssetId int32       `json:"assetId"`
	Subtype HeapSubtype `json:"subtype"`
	// Timestamp of the latest data change
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// Asset payload
	Data map[string]interface{} `json:"data"`
}

// NewHeap instantiates a new Heap object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHeap(assetId int32, subtype HeapSubtype, data map[string]interface{}) *Heap {
	this := Heap{}
	this.AssetId = assetId
	this.Subtype = subtype
	this.Data = data
	return &this
}

// NewHeapWithDefaults instantiates a new Heap object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHeapWithDefaults() *Heap {
	this := Heap{}
	var subtype HeapSubtype = INPUT
	this.Subtype = subtype
	return &this
}

// GetAssetId returns the AssetId field value
func (o *Heap) GetAssetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.AssetId
}

// GetAssetIdOk returns a tuple with the AssetId field value
// and a boolean to check if the value has been set.
func (o *Heap) GetAssetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AssetId, true
}

// SetAssetId sets field value
func (o *Heap) SetAssetId(v int32) {
	o.AssetId = v
}

// GetSubtype returns the Subtype field value
func (o *Heap) GetSubtype() HeapSubtype {
	if o == nil {
		var ret HeapSubtype
		return ret
	}

	return o.Subtype
}

// GetSubtypeOk returns a tuple with the Subtype field value
// and a boolean to check if the value has been set.
func (o *Heap) GetSubtypeOk() (*HeapSubtype, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Subtype, true
}

// SetSubtype sets field value
func (o *Heap) SetSubtype(v HeapSubtype) {
	o.Subtype = v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise.
func (o *Heap) GetTimestamp() time.Time {
	if o == nil || o.Timestamp == nil {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Heap) GetTimestampOk() (*time.Time, bool) {
	if o == nil || o.Timestamp == nil {
		return nil, false
	}
	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *Heap) HasTimestamp() bool {
	if o != nil && o.Timestamp != nil {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *Heap) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}

// GetData returns the Data field value
func (o *Heap) GetData() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *Heap) GetDataOk() (map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *Heap) SetData(v map[string]interface{}) {
	o.Data = v
}

func (o Heap) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["assetId"] = o.AssetId
	}
	if true {
		toSerialize["subtype"] = o.Subtype
	}
	if o.Timestamp != nil {
		toSerialize["timestamp"] = o.Timestamp
	}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableHeap struct {
	value *Heap
	isSet bool
}

func (v NullableHeap) Get() *Heap {
	return v.value
}

func (v *NullableHeap) Set(val *Heap) {
	v.value = val
	v.isSet = true
}

func (v NullableHeap) IsSet() bool {
	return v.isSet
}

func (v *NullableHeap) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHeap(val *Heap) *NullableHeap {
	return &NullableHeap{value: val, isSet: true}
}

func (v NullableHeap) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHeap) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
