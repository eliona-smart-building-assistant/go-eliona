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

package asset

import (
	"fmt"
	"reflect"
	"strings"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v2/tools"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
)

// UpsertData inserts or updates the given asset data. If the data with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertData(data api.Data) error {
	_, err := client.NewClient().DataApi.
		PutData(client.AuthenticationContext()).
		Data(data).
		Execute()
	tools.LogError(err)
	return err
}

// UpsertDataIfAssetExists upserts the data if the eliona id exists. Otherwise, the upsert is ignored.
func UpsertDataIfAssetExists(data api.Data) error {
	exists, err := ExistAsset(data.AssetId)
	if err != nil {
		return err
	}
	if exists {
		return UpsertData(data)
	}
	return nil
}

type Data struct {
	AssetId   int32
	Timestamp api.NullableTime
	Data      any
}

// UpsertAssetDataIfAssetExists upserts the data in any struct having `eliona` field tags.
// If the eliona ID does not exist, the upsert is ignored.
func UpsertAssetDataIfAssetExists(data Data) error {
	a, err := getAsset(data.AssetId)
	if err != nil {
		return fmt.Errorf("getting asset id %v: %v", data.AssetId, err)
	}
	if a == nil {
		return nil
	}
	asset := *a

	subtypes := splitBySubtype(data)
	for subtype, subData := range subtypes {
		UpsertData(api.Data{
			AssetId:       data.AssetId,
			Subtype:       subtype,
			Timestamp:     data.Timestamp,
			Data:          subData,
			AssetTypeName: *api.NewNullableString(&asset.AssetType),
		})
	}

	return nil
}

func splitBySubtype(data any) map[api.DataSubtype]map[string]interface{} {
	value := reflect.ValueOf(data)
	valueType := reflect.TypeOf(data)

	result := make(map[api.DataSubtype]map[string]interface{})

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i).Interface()

		elionaTag := field.Tag.Get("eliona")
		// Omit the attribute properties
		attributeName := strings.Split(elionaTag, ",")[0]

		subtype := api.DataSubtype(field.Tag.Get("subtype"))

		if subtype == api.SUBTYPE_OUTPUT {
			continue
		}

		if _, ok := result[subtype]; !ok {
			result[subtype] = make(map[string]interface{})
		}
		result[subtype][attributeName] = fieldValue
	}

	return result
}
