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

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

type AssetWithParentReferences interface {
	GetName() string
	GetDescription() string
	GetAssetType() string
	GetGAI() string
	GetLocationalParentGAI() string
	GetFunctionalParentGAI() string

	SetAssetID(assetID int32, projectID string) error
}

// CreateAssetsBulk creates assets within the asset structure provided.
// It ensures that parent assets are created before their children by sorting the assets accordingly.
func CreateAssetsBulk(assets []AssetWithParentReferences, projectId string) (createdCnt int, err error) {
	return createAssets(assets, projectId)
}

func createAssets(assets []AssetWithParentReferences, projectId string) (createdCnt int, err error) {
	sortedAssets, err := sortAssetsByDependencies(assets)
	if err != nil {
		return 0, fmt.Errorf("sorting assets: %v", err)
	}

	var apiAssets []api.Asset
	for _, a := range sortedAssets {
		apiAssets = append(apiAssets, assetToAPIAsset(a, projectId))
	}
	result, err := UpsertAssetsBulkGAI(apiAssets)
	if err != nil {
		return 0, fmt.Errorf("upserting bulk assets: %v", err)
	}

	resultsMap := make(map[string]api.Asset)
	for _, a := range result {
		resultsMap[a.GlobalAssetIdentifier] = a
	}
	for _, a := range sortedAssets {
		updatedAPIAsset, ok := resultsMap[a.GetGAI()]
		if !ok {
			return 0, fmt.Errorf("should not happen: did not find GAI '%v' in updated api assets", a.GetGAI())
		}
		if !updatedAPIAsset.Id.IsSet() || *updatedAPIAsset.Id.Get() == 0 {
			return 0, fmt.Errorf("GAI '%v' has zero value AssetID in response from APIv2", a.GetGAI())
		}
		if err := a.SetAssetID(updatedAPIAsset.GetId(), projectId); err != nil {
			return 0, fmt.Errorf("setting asset ID: %v", err)
		}
	}

	return len(result), nil
}

// sortAssetsByDependencies sorts the assets to ensure that parent assets are created before their children.
// It performs a topological sort based on the parent GAIs.
func sortAssetsByDependencies(assets []AssetWithParentReferences) ([]AssetWithParentReferences, error) {
	assetMap := make(map[string]AssetWithParentReferences)
	for _, asset := range assets {
		assetMap[asset.GetGAI()] = asset
	}

	var sortedAssets []AssetWithParentReferences
	visited := make(map[string]bool)
	tempMark := make(map[string]bool)

	var visit func(string) error
	visit = func(gai string) error {
		if tempMark[gai] {
			return fmt.Errorf("circular dependency detected at asset with GAI '%s'", gai)
		}
		if visited[gai] {
			return nil
		}

		tempMark[gai] = true
		defer func() { tempMark[gai] = false }()

		asset, exists := assetMap[gai]
		if !exists {
			return fmt.Errorf("asset with GAI '%s' not found", gai)
		}

		locParentGAI := asset.GetLocationalParentGAI()
		if locParentGAI != "" {
			if _, parentExists := assetMap[locParentGAI]; parentExists {
				if err := visit(locParentGAI); err != nil {
					return err
				}
			}
		}

		funcParentGAI := asset.GetFunctionalParentGAI()
		if funcParentGAI != "" {
			if _, parentExists := assetMap[funcParentGAI]; parentExists {
				if err := visit(funcParentGAI); err != nil {
					return err
				}
			}
		}

		visited[gai] = true
		sortedAssets = append(sortedAssets, asset)
		return nil
	}

	for gai := range assetMap {
		if !visited[gai] {
			if err := visit(gai); err != nil {
				return nil, err
			}
		}
	}

	return sortedAssets, nil
}

func assetToAPIAsset(ast AssetWithParentReferences, projectId string) api.Asset {
	return api.Asset{
		ProjectId:                  projectId,
		GlobalAssetIdentifier:      ast.GetGAI(),
		Name:                       *api.NewNullableString(common.Ptr(ast.GetName())),
		AssetType:                  ast.GetAssetType(),
		Description:                *api.NewNullableString(common.Ptr(ast.GetDescription())),
		ParentFunctionalIdentifier: *api.NewNullableString(common.Ptr(ast.GetFunctionalParentGAI())),
		ParentLocationalIdentifier: *api.NewNullableString(common.Ptr(ast.GetLocationalParentGAI())),
	}
}
