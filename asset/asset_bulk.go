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
	"fmt"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

type AssetLikeWithParentReferences interface {
	GetName() string
	GetDescription() string
	GetAssetType() string
	GetGAI() string
	GetLocationalParentGAI() string
	GetFunctionalParentGAI() string
	GetSiteID() string

	SetAssetID(assetID int32) error
}

// CreateAssetsBulk creates assets within the asset structure provided.
// It ensures that parent assets are created before their children by sorting the assets accordingly.
func CreateAssetsBulk(apiEndpoint string, apiKey string, assetLikes []AssetLikeWithParentReferences) (createdCnt int, err error) {
	return createAssets(apiEndpoint, apiKey, assetLikes)
}

func createAssets(apiEndpoint string, apiKey string, assetLikes []AssetLikeWithParentReferences) (createdCnt int, err error) {
	sortedAssetLikes, err := sortAssetLikesByDependencies(assetLikes)
	if err != nil {
		return 0, fmt.Errorf("sorting assetLikes: %v", err)
	}

	var apiAssets []api.Asset
	for _, assetLike := range sortedAssetLikes {
		apiAssets = append(apiAssets, assetLikeToApiAsset(assetLike))
	}
	result, err := UpsertAssetsBulkGAI(apiEndpoint, apiKey, apiAssets)
	if err != nil {
		return 0, fmt.Errorf("upserting bulk assetLikes: %v", err)
	}

	resultsMap := make(map[string]api.Asset)
	for _, a := range result {
		resultsMap[a.GlobalAssetIdentifier] = a
	}
	for _, assetLike := range sortedAssetLikes {
		updatedApiAsset, ok := resultsMap[assetLike.GetGAI()]
		if !ok {
			return 0, fmt.Errorf("should not happen: did not find GAI '%v' in updated api assetLike", assetLike.GetGAI())
		}
		if !updatedApiAsset.Id.IsSet() || *updatedApiAsset.Id.Get() == 0 {
			return 0, fmt.Errorf("GAI '%v' has zero value AssetID in response from APIv2", assetLike.GetGAI())
		}
		if err := assetLike.SetAssetID(updatedApiAsset.GetId()); err != nil {
			return 0, fmt.Errorf("setting asset ID: %v", err)
		}
	}

	return len(result), nil
}

// sortAssetLikesByDependencies sorts the assets to ensure that parent assets are created before their children.
// It performs a topological sort based on the parent GAIs.
func sortAssetLikesByDependencies(assetLikes []AssetLikeWithParentReferences) ([]AssetLikeWithParentReferences, error) {
	assetLikeMap := make(map[string]AssetLikeWithParentReferences)
	for _, assetLike := range assetLikes {
		assetLikeMap[assetLike.GetGAI()] = assetLike
	}

	var sortedAssetLikes []AssetLikeWithParentReferences
	visited := make(map[string]bool)
	tempMark := make(map[string]bool)

	var visit func(string) error
	visit = func(gai string) error {
		if tempMark[gai] {
			return fmt.Errorf("circular dependency detected at assetLike with GAI '%s'", gai)
		}
		if visited[gai] {
			return nil
		}

		tempMark[gai] = true
		defer func() { tempMark[gai] = false }()

		asset, exists := assetLikeMap[gai]
		if !exists {
			return fmt.Errorf("assetLike with GAI '%s' not found", gai)
		}

		locParentGAI := asset.GetLocationalParentGAI()
		if locParentGAI != "" {
			if _, parentExists := assetLikeMap[locParentGAI]; parentExists {
				if err := visit(locParentGAI); err != nil {
					return err
				}
			}
		}

		funcParentGAI := asset.GetFunctionalParentGAI()
		if funcParentGAI != "" {
			if _, parentExists := assetLikeMap[funcParentGAI]; parentExists {
				if err := visit(funcParentGAI); err != nil {
					return err
				}
			}
		}

		visited[gai] = true
		sortedAssetLikes = append(sortedAssetLikes, asset)
		return nil
	}

	for gai := range assetLikeMap {
		if !visited[gai] {
			if err := visit(gai); err != nil {
				return nil, err
			}
		}
	}

	return sortedAssetLikes, nil
}

func assetLikeToApiAsset(assetLike AssetLikeWithParentReferences) api.Asset {
	return api.Asset{
		SiteId:                     *api.NewNullableString(common.Ptr(assetLike.GetSiteID())),
		GlobalAssetIdentifier:      assetLike.GetGAI(),
		Name:                       *api.NewNullableString(common.Ptr(assetLike.GetName())),
		AssetType:                  assetLike.GetAssetType(),
		Description:                *api.NewNullableString(common.Ptr(assetLike.GetDescription())),
		ParentFunctionalIdentifier: *api.NewNullableString(common.Ptr(assetLike.GetFunctionalParentGAI())),
		ParentLocationalIdentifier: *api.NewNullableString(common.Ptr(assetLike.GetLocationalParentGAI())),
	}
}
