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
func CreateAssetsBulk(assets []AssetWithParentReferences, projectId string) (createdCnt int, err error) {
	return createAssets(assets, projectId)
}

func createAssets(assets []AssetWithParentReferences, projectId string) (createdCnt int, err error) {
	var apiAssets []api.Asset
	for _, a := range assets {
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
	for _, a := range assets {
		updatedAPIAsset, ok := resultsMap[a.GetGAI()]
		if !ok {
			return 0, fmt.Errorf("should not happen: did not find GAI '%v' in updated api assets", a.GetGAI())
		}
		if !updatedAPIAsset.Id.IsSet() || *updatedAPIAsset.Id.Get() == 0 {
			return 0, fmt.Errorf("GAI '%v' has zero valueAssetID in response from APIv2", a.GetGAI())
		}
		if err := a.SetAssetID(updatedAPIAsset.GetId(), projectId); err != nil {
			return 0, fmt.Errorf("setting asset ID: %v", err)
		}
	}

	return len(result), nil
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
