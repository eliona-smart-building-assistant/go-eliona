//  This file is part of the eliona project.
//  Copyright © 2024 LEICOM iTEC AG. All Rights Reserved.
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
	"net/http"
	"path/filepath"

	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v3/tools"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
	"github.com/eliona-smart-building-assistant/go-eliona/v2/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

var ErrNotFound = errors.New("not found")

// UpsertAssetType insert or, when already exist, updates an asset type
func UpsertAssetType(apiEndpoint string, apiKey string, assetType api.AssetType) error {
	_, _, err := client.NewClient(apiEndpoint).AssetTypesAPI.
		PutAssetType(client.AuthenticationContext(apiKey)).
		Expansions([]string{"AssetType.attributes"}). // take values of attributes also
		AssetType(assetType).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting asset type %v: %w", assetType.Name, err))
	}
	return err
}

func getAsset(apiEndpoint string, apiKey string, assetId int32) (*api.Asset, error) {
	asset, res, err := client.NewClient(apiEndpoint).AssetsAPI.
		GetAssetById(client.AuthenticationContext(apiKey), assetId).
		Execute()
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	return asset, err
}

// ExistAsset returns true, if the given asset id exists in eliona
func ExistAsset(apiEndpoint string, apiKey string, assetId int32) (bool, error) {
	_, err := getAsset(apiEndpoint, apiKey, assetId)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		tools.LogError(fmt.Errorf("checking if asset %v exists: %w", assetId, err))
	}
	return true, err
}

// UpsertAsset inserts or updates an asset and returns the id
func UpsertAsset(apiEndpoint string, apiKey string, asset api.Asset) (*int32, error) {
	upsertedAsset, _, err := client.NewClient(apiEndpoint).AssetsAPI.
		PutAsset(client.AuthenticationContext(apiKey)).
		Asset(asset).Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting asset %v: %w", asset.Name, err))
	}
	if err != nil {
		return nil, err
	}
	return upsertedAsset.Id.Get(), nil
}

// UpsertAssetsBulkGAI inserts or updates all given assets and returns them
// including asset IDs. Relations are identified by GAI + SiteID.
// Assets must be ordered - parents must come before their children, otherwise
// relations will not be created.
func UpsertAssetsBulkGAI(apiEndpoint string, apiKey string, assets []api.Asset) ([]api.Asset, error) {
	upsertedAssets, _, err := client.NewClient(apiEndpoint).AssetsAPI.
		PutBulkAssets(client.AuthenticationContext(apiKey)).
		Asset(assets).
		IdentifyBy(string(api.ASSET_IDENTIFY_BY_GAI_SITE_ID)).
		Execute()
	if err != nil {
		e := fmt.Errorf("upserting assets: %w", err)
		tools.LogError(e)
		return nil, e
	}
	return upsertedAssets, nil
}

// UpsertAssetTypeAttribute insert or updates an asset and returns the id
func UpsertAssetTypeAttribute(apiEndpoint string, apiKey string, attribute api.AssetTypeAttribute) error {
	_, _, err := client.NewClient(apiEndpoint).AssetTypesAPI.
		PutAssetTypeAttribute(client.AuthenticationContext(apiKey), *attribute.AssetTypeName.Get()).
		AssetTypeAttribute(attribute).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("upserting asset type attribute %v: %w", attribute.Name, err))
	}
	return err
}

// InitAssetType inserts or updates the given asset type.
func InitAssetType(apiEndpoint string, apiKey string, assetType api.AssetType) func(db.Connection) error {
	return func(db.Connection) error {
		return UpsertAssetType(apiEndpoint, apiKey, assetType)
	}
}

func initAssetTypeFile(apiEndpoint string, apiKey string, path string) error {
	assetType, err := common.UnmarshalFile[api.AssetType](path)
	if err != nil {
		return fmt.Errorf("unmarshalling file %s: %v", path, err)
	}
	return UpsertAssetType(apiEndpoint, apiKey, assetType)
}

// InitAssetTypeFile inserts or updates the asset type build from the content of the given file.
func InitAssetTypeFile(apiEndpoint string, apiKey string, path string) func(db.Connection) error {
	return func(db.Connection) error {
		return initAssetTypeFile(apiEndpoint, apiKey, path)
	}
}

func InitAssetTypeFiles(apiEndpoint string, apiKey string, pattern string) func(db.Connection) error {
	return func(db.Connection) error {
		paths, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("glob file pattern %s: %v", pattern, err)
		}
		for _, path := range paths {
			err := initAssetTypeFile(apiEndpoint, apiKey, path)
			if err != nil {
				return fmt.Errorf("initializing asset type %s: %v", path, err)
			}
		}
		return nil
	}
}

type SubType string

const (
	Status   SubType = "status"
	Info     SubType = "info"
	Input    SubType = "input"
	Output   SubType = "output"
	Property SubType = "property"
)
