//  This file is part of the eliona project.
//  Copyright © 2025 LEICOM iTEC AG. All Rights Reserved.
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

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v3/tools"
	"github.com/eliona-smart-building-assistant/go-eliona/v2/client"
)

// UpsertData inserts or updates the given asset data. If the data with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertData(apiEndpoint string, apiKey string, data api.Data) error {
	_, err := client.NewClient(apiEndpoint).DataAPI.
		PutData(client.AuthenticationContext(apiKey)).
		Data(data).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting data for asset %v: %w", data.AssetId, err))
	}
	return err
}

// UpsertDataBulk inserts or updates the given asset data. If the data with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertDataBulk(apiEndpoint string, apiKey string, datas []api.Data) error {
	_, err := client.NewClient(apiEndpoint).DataAPI.
		PutBulkData(client.AuthenticationContext(apiKey)).
		Data(datas).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting data bulk: %w", err))
	}
	return err
}

// UpsertDataIfAssetExists upserts the data if the eliona id exists. Otherwise, the upsert is ignored.
func UpsertDataIfAssetExists(apiEndpoint string, apiKey string, data api.Data) error {
	exists, err := ExistAsset(apiEndpoint, apiKey, data.AssetId)
	if err != nil {
		return fmt.Errorf("checking if asset %v exists: %w", data.AssetId, err)
	}
	if exists {
		return UpsertData(apiEndpoint, apiKey, data)
	}
	return nil
}

// UpsertDataBulkIfAssetExists upserts the data if the eliona id exists. Otherwise, the upsert is ignored.
func UpsertDataBulkIfAssetExists(apiEndpoint string, apiKey string, datas []api.Data) error {
	upsertDatas := make([]api.Data, 0, len(datas))
	for _, data := range datas {
		exists, err := ExistAsset(apiEndpoint, apiKey, data.AssetId)
		if err != nil {
			return fmt.Errorf("checking if asset %v exists: %w", data.AssetId, err)
		}
		if exists {
			upsertDatas = append(upsertDatas, data)
		}
	}
	return UpsertDataBulk(apiEndpoint, apiKey, upsertDatas)
}

type Data struct {
	AssetId         int32
	Timestamp       api.NullableTime
	ClientReference string
	Data            any
}

// UpsertAssetDataIfAssetExists upserts the data in any struct having `eliona` field tags.
// If the eliona ID does not exist, the upsert is ignored.
func UpsertAssetDataIfAssetExists(apiEndpoint string, apiKey string, data Data) error {
	subtypes := SplitBySubtype(data.Data)
	for subtype, subData := range subtypes {
		if err := UpsertData(apiEndpoint, apiKey, api.Data{
			AssetId:         data.AssetId,
			Subtype:         subtype,
			Timestamp:       data.Timestamp,
			Data:            subData,
			ClientReference: *api.NewNullableString(&data.ClientReference),
		}); err != nil {
			return fmt.Errorf("upserting data for subtype %s: %v", subtype, err)
		}
	}

	return nil
}

func SplitBySubtype(data any) map[api.DataSubtype]map[string]interface{} {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	valueType := reflect.TypeOf(data)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	result := make(map[api.DataSubtype]map[string]interface{})

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)

		if field.PkgPath != "" {
			// Skip unexported fields.
			continue
		}
		fieldValue := value.Field(i).Interface()

		tag, ok := ParseElionaTag(field)
		if !ok {
			// Skip fields without tag.
			continue
		}

		if tag.Subtype == "" {
			continue
		}

		if _, ok := result[tag.Subtype]; !ok {
			result[tag.Subtype] = make(map[string]interface{})
		}
		result[tag.Subtype][tag.AttributeName] = fieldValue
	}

	return result
}

func GetData(apiEndpoint string, apiKey string, assetID int32, subtype string) ([]api.Data, error) {
	data, _, err := client.NewClient(apiEndpoint).DataAPI.
		GetData(client.AuthenticationContext(apiKey)).
		AssetId(assetID).
		DataSubtype(subtype).
		Execute()
	if err != nil {
		err = fmt.Errorf("getting data for asset %v subtype %v: %w", assetID, subtype, err)
		tools.LogError(err)
		return nil, err
	}
	return data, nil
}
