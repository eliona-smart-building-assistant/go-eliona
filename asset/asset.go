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
	"context"
	"github.com/eliona-smart-building-assistant/go-eliona/api"
)

// UpsertAssetType insert or, when already exist, updates an asset type
func UpsertAssetType(assetType api.AssetType) error {
	_, err := api.NewClient().AssetTypeApi.PostAssetType(context.Background()).AssetType(assetType).Execute()
	return err
}

// ExistAsset returns true, if the given asset id exists in eliona
func ExistAsset(assetId int32) (bool, error) {
	asset, _, err := api.NewClient().AssetApi.GetAssetById(context.Background(), assetId).Execute()
	return asset != nil, err
}

// UpsertAsset insert or updates an assetdb and returns the id
func UpsertAsset(asset api.Asset) (*int32, error) {
	upsertedAsset, _, err := api.NewClient().AssetApi.PostAsset(context.Background()).Asset(asset).Execute()
	return upsertedAsset.Id, err
}

// UpsertAssetTypeAttribute insert or updates an assetdb and returns the id
func UpsertAssetTypeAttribute(attribute api.Attribute) error {
	_, err := api.NewClient().AssetTypeApi.PostAssetTypeAttribute(context.Background()).Attribute(attribute).Execute()
	return err
}
