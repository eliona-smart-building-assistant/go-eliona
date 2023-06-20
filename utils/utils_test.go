//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package utils

import (
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDeviceInfo struct {
	Model        string `json:"model" eliona:"model,filterable"`
	ID           string `json:"id" eliona:"id,filterable"`
	BatteryLevel int    `json:"batteryLevel" eliona:"battery_level" subtype:"status"`
	Product      string `json:"product" eliona:"product,filterable"`
	Name         string `json:"name" eliona:"name,filterable"`
	Mac          string `json:"mac" eliona:"mac,filterable"`
	Firmware     string `json:"firmware" eliona:"firm_ware,filterable"`
}

func TestParseElionaTag(t *testing.T) {
	input := testDeviceInfo{}
	inputType := reflect.TypeOf(input)

	elionaTag, err := parseElionaTag(inputType.Field(0))
	assert.NoError(t, err)
	assert.Equal(t, true, elionaTag.Filterable)
	assert.Equal(t, "model", elionaTag.ParamName)
	assert.Equal(t, asset.SubType(""), elionaTag.SubType)

	elionaTag, err = parseElionaTag(inputType.Field(2))
	assert.NoError(t, err)
	assert.Equal(t, false, elionaTag.Filterable)
	assert.Equal(t, "battery_level", elionaTag.ParamName)
	assert.Equal(t, asset.Status, elionaTag.SubType)
}

func TestStructToMap(t *testing.T) {
	input := testDeviceInfo{Name: "test", Firmware: "xyz"}
	output, err := StructToMap(input)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(output))
	assert.Equal(t, "test", output["name"])
	assert.Equal(t, "xyz", output["firm_ware"])
}
