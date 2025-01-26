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
	"time"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
)

type Root interface {
	LocationalNode
	FunctionalNode
}

type LocationalNode interface {
	Asset
	GetLocationalChildren() []LocationalNode
}

type FunctionalNode interface {
	Asset
	GetFunctionalChildren() []FunctionalNode
}

type Asset interface {
	GetName() string
	GetDescription() string
	GetAssetType() string
	GetGAI() string

	GetAssetID(projectID string) (*int32, error)
	SetAssetID(assetID int32, projectID string) error
}

// CreateAssets creates assets within the asset structure using bulk creation.
// It skips assets that have assetID != nil but keeps their children.
// NOTE: Prefer using CreateAssetsBulk method, avoiding complexity and too much
// juggling of asset structure.
func CreateAssets(root Root, projectId string) (createdCnt int, err error) {
	var assetsToCreate []AssetWithParentReferences

	err = collectAssetsToCreate(root, "", "", &assetsToCreate, projectId, map[string]bool{})
	if err != nil {
		return 0, fmt.Errorf("collecting assets to create: %v", err)
	}

	createdCnt, err = createAssets(assetsToCreate, projectId)
	if err != nil {
		return 0, fmt.Errorf("creating assets: %v", err)
	}

	return createdCnt, nil
}

// CreateAssetsAndUpsertData creates assets within the asset structure and upserts their data.
// It skips assets that have assetID != nil but keeps their children.
// It uses the createAssets() function to bulk create assets.
// NOTE: Prefer using CreateAssetsBulk method, avoiding complexity and too much
// juggling of asset structure.
func CreateAssetsAndUpsertData(root Root, projectId string, ts *time.Time, clientReference *string) (createdCnt int, err error) {
	var assetsToCreate []AssetWithParentReferences
	var dataToUpsert []Data

	err = collectAssetsAndDataToCreate(root, "", "", &assetsToCreate, &dataToUpsert, projectId, ts, clientReference, map[string]bool{})
	if err != nil {
		return 0, fmt.Errorf("collecting assets and data to create: %v", err)
	}

	createdCnt, err = createAssets(assetsToCreate, projectId)
	if err != nil {
		return 0, fmt.Errorf("creating assets: %v", err)
	}

	// Upsert data for all assets (including those that already existed)
	for _, data := range dataToUpsert {
		err := UpsertAssetDataIfAssetExists(data)
		if err != nil {
			return createdCnt, fmt.Errorf("upserting data: %v", err)
		}
	}

	return createdCnt, nil
}

// collectAssetsToCreate traverses the asset tree and collects assets that need to be created.
// Note: The parent asset is always added to the assets slice before its children.
func collectAssetsToCreate(node Asset, locationalParentGAI, functionalParentGAI string, assets *[]AssetWithParentReferences, projectId string, visited map[string]bool) error {
	if visited[node.GetGAI()] {
		// Update existing relationships instead of skipping entirely
		for _, asset := range *assets {
			if asset.GetGAI() == node.GetGAI() {
				// Update locational and functional parent relationships
				if locationalParentGAI != "" && asset.GetLocationalParentGAI() == "" {
					asset.(*AssetToCreate).locationalParentGAI = locationalParentGAI
				}
				if functionalParentGAI != "" && asset.GetFunctionalParentGAI() == "" {
					asset.(*AssetToCreate).functionalParentGAI = functionalParentGAI
				}
				break
			}
		}
		return nil
	}
	visited[node.GetGAI()] = true

	assetID, err := node.GetAssetID(projectId)
	if err != nil {
		return fmt.Errorf("getting asset ID: %v", err)
	}

	if assetID == nil {
		a := &AssetToCreate{
			node:                node,
			locationalParentGAI: locationalParentGAI,
			functionalParentGAI: functionalParentGAI,
		}
		*assets = append(*assets, a)
	}

	if ln, ok := node.(LocationalNode); ok {
		for _, child := range ln.GetLocationalChildren() {
			if child == nil {
				continue
			}
			err := collectAssetsToCreate(child, node.GetGAI(), functionalParentGAI, assets, projectId, visited)
			if err != nil {
				return err
			}
		}
	}

	if fn, ok := node.(FunctionalNode); ok {
		for _, child := range fn.GetFunctionalChildren() {
			if child == nil {
				continue
			}
			err := collectAssetsToCreate(child, locationalParentGAI, node.GetGAI(), assets, projectId, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func collectAssetsAndDataToCreate(node Asset, locationalParentGAI, functionalParentGAI string, assets *[]AssetWithParentReferences, data *[]Data, projectId string, ts *time.Time, clientReference *string, visited map[string]bool) error {
	if visited[node.GetGAI()] {
		// Update existing relationships instead of skipping entirely
		for _, asset := range *assets {
			if asset.GetGAI() == node.GetGAI() {
				// Update locational and functional parent relationships
				if locationalParentGAI != "" && asset.GetLocationalParentGAI() == "" {
					asset.(*AssetToCreate).locationalParentGAI = locationalParentGAI
				}
				if functionalParentGAI != "" && asset.GetFunctionalParentGAI() == "" {
					asset.(*AssetToCreate).functionalParentGAI = functionalParentGAI
				}
				break
			}
		}
		return nil
	}
	visited[node.GetGAI()] = true

	assetID, err := node.GetAssetID(projectId)
	if err != nil {
		return fmt.Errorf("getting asset ID: %v", err)
	}

	if assetID == nil {
		a := &AssetToCreate{
			node:                node,
			locationalParentGAI: locationalParentGAI,
			functionalParentGAI: functionalParentGAI,
		}
		*assets = append(*assets, a)
	} else {
		// Prepare data to upsert for existing assets
		dataToUpsert := Data{
			AssetId:         *assetID,
			Timestamp:       *api.NewNullableTime(ts),
			ClientReference: *clientReference,
			Data:            node,
		}
		*data = append(*data, dataToUpsert)
	}

	if ln, ok := node.(LocationalNode); ok {
		for _, child := range ln.GetLocationalChildren() {
			if child == nil {
				continue
			}
			err := collectAssetsAndDataToCreate(child, node.GetGAI(), functionalParentGAI, assets, data, projectId, ts, clientReference, visited)
			if err != nil {
				return err
			}
		}
	}

	if fn, ok := node.(FunctionalNode); ok {
		for _, child := range fn.GetFunctionalChildren() {
			if child == nil {
				continue
			}
			err := collectAssetsAndDataToCreate(child, locationalParentGAI, node.GetGAI(), assets, data, projectId, ts, clientReference, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AssetToCreate implements AssetWithParentReferences and holds asset information for creation.
type AssetToCreate struct {
	node                Asset
	locationalParentGAI string
	functionalParentGAI string
}

func (a *AssetToCreate) GetName() string {
	return a.node.GetName()
}

func (a *AssetToCreate) GetDescription() string {
	return a.node.GetDescription()
}

func (a *AssetToCreate) GetAssetType() string {
	return a.node.GetAssetType()
}

func (a *AssetToCreate) GetGAI() string {
	return a.node.GetGAI()
}

func (a *AssetToCreate) GetLocationalParentGAI() string {
	return a.locationalParentGAI
}

func (a *AssetToCreate) GetFunctionalParentGAI() string {
	return a.functionalParentGAI
}

func (a *AssetToCreate) SetAssetID(assetID int32, projectID string) error {
	return a.node.SetAssetID(assetID, projectID)
}
