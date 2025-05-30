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
	"errors"
	"fmt"
	"reflect"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v2/tools"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
)

// UpsertData inserts or updates the given asset data. If the data with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertData(data api.Data) error {
	_, err := client.NewClient().DataAPI.
		PutData(client.AuthenticationContext()).
		Data(data).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting data for asset %v: %w", data.AssetId, err))
	}
	return err
}

// UpsertDataBulk inserts or updates the given asset data. If the data with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertDataBulk(datas []api.Data) error {
	_, err := client.NewClient().DataAPI.
		PutBulkData(client.AuthenticationContext()).
		Data(datas).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting data bulk: %w", err))
	}
	return err
}

// UpsertDataIfAssetExists upserts the data if the eliona id exists. Otherwise, the upsert is ignored.
func UpsertDataIfAssetExists(data api.Data) error {
	exists, err := ExistAsset(data.AssetId)
	if err != nil {
		return fmt.Errorf("checking if asset %v exists: %w", data.AssetId, err)
	}
	if exists {
		return UpsertData(data)
	}
	return nil
}

// UpsertDataBulkIfAssetExists upserts the data if the eliona id exists. Otherwise, the upsert is ignored.
func UpsertDataBulkIfAssetExists(datas []api.Data) error {
	upsertDatas := make([]api.Data, 0, len(datas))
	for _, data := range datas {
		exists, err := ExistAsset(data.AssetId)
		if err != nil {
			return fmt.Errorf("checking if asset %v exists: %w", data.AssetId, err)
		}
		if exists {
			upsertDatas = append(upsertDatas, data)
		}
	}
	return UpsertDataBulk(upsertDatas)
}

type Data struct {
	AssetId         int32
	Timestamp       api.NullableTime
	ClientReference string
	Data            any
}

// UpsertAssetDataIfAssetExists upserts the data in any struct having `eliona` field tags.
// If the eliona ID does not exist, the upsert is ignored.
func UpsertAssetDataIfAssetExists(data Data) error {
	a, err := getAsset(data.AssetId)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("getting asset id %v: %v", data.AssetId, err)
	}
	if a == nil {
		return fmt.Errorf("shouldn't happen: asset with id %v is nil", data.AssetId)
	}
	asset := *a

	subtypes := SplitBySubtype(data.Data)
	for subtype, subData := range subtypes {
		if err := UpsertData(api.Data{
			AssetId:         data.AssetId,
			Subtype:         subtype,
			Timestamp:       data.Timestamp,
			Data:            subData,
			AssetTypeName:   *api.NewNullableString(&asset.AssetType),
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

func GetData(assetID int32, subtype string) ([]api.Data, error) {
	data, _, err := client.NewClient().DataAPI.
		GetData(client.AuthenticationContext()).
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
