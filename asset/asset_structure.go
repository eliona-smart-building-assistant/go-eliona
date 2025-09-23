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

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
)

type Root interface {
	LocationalNode
	FunctionalNode
}

type LocationalNode interface {
	AssetLike
	GetLocationalChildren() []LocationalNode
}

type FunctionalNode interface {
	AssetLike
	GetFunctionalChildren() []FunctionalNode
}

type AssetLike interface {
	GetName() string
	GetDescription() string
	GetAssetType() string
	GetGAI() string
	GetSiteID() string

	GetAssetID() (*int32, error)
	SetAssetID(assetID int32) error
}

// CreateAssets creates assets within the asset like structure using bulk creation.
// It skips asset likes that have assetID != nil but keeps their children.
// NOTE: Prefer using CreateAssetsBulk method, avoiding complexity and too much
// juggling of asset structure.
func CreateAssets(apiEndpoint string, apiKey string, root Root) (createdCnt int, err error) {
	var assetLikesToCreate []AssetLikeWithParentReferences

	err = collectAssetLikesToCreate(root, "", "", &assetLikesToCreate, map[string]bool{})
	if err != nil {
		return 0, fmt.Errorf("collecting asset likes to create: %v", err)
	}

	createdCnt, err = createAssets(apiEndpoint, apiKey, assetLikesToCreate)
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
func CreateAssetsAndUpsertData(apiEndpoint string, apiKey string, root Root, ts *time.Time, clientReference *string) (createdCnt int, err error) {
	var assetLikesToCreate []AssetLikeWithParentReferences
	var dataToUpsert []Data

	err = collectAssetLikesAndDataToCreate(root, "", "", &assetLikesToCreate, &dataToUpsert, ts, clientReference, map[string]bool{})
	if err != nil {
		return 0, fmt.Errorf("collecting asset likes and data to create: %v", err)
	}

	createdCnt, err = createAssets(apiEndpoint, apiKey, assetLikesToCreate)
	if err != nil {
		return 0, fmt.Errorf("creating assets: %v", err)
	}

	// Upsert data for all assets (including those that already existed)
	for _, data := range dataToUpsert {
		err := UpsertAssetDataIfAssetExists(apiEndpoint, apiKey, data)
		if err != nil {
			return createdCnt, fmt.Errorf("upserting data: %v", err)
		}
	}

	return createdCnt, nil
}

// collectAssetLikesToCreate traverses the asset like tree and collects asset likes that need to be created.
// Note: The parent asset like is always added to the asset likes slice before its children.
func collectAssetLikesToCreate(node AssetLike, locationalParentGAI, functionalParentGAI string, assets *[]AssetLikeWithParentReferences, visited map[string]bool) error {
	if visited[node.GetGAI()] {
		// Update existing relationships instead of skipping entirely
		for _, asset := range *assets {
			if asset.GetGAI() == node.GetGAI() {
				// Update locational and functional parent relationships
				if locationalParentGAI != "" && asset.GetLocationalParentGAI() == "" {
					asset.(*AssetLikeToCreate).locationalParentGAI = locationalParentGAI
				}
				if functionalParentGAI != "" && asset.GetFunctionalParentGAI() == "" {
					asset.(*AssetLikeToCreate).functionalParentGAI = functionalParentGAI
				}
				break
			}
		}
		return nil
	}
	visited[node.GetGAI()] = true

	assetID, err := node.GetAssetID()
	if err != nil {
		return fmt.Errorf("getting asset ID: %v", err)
	}

	if assetID == nil {
		a := &AssetLikeToCreate{
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
			err := collectAssetLikesToCreate(child, node.GetGAI(), functionalParentGAI, assets, visited)
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
			err := collectAssetLikesToCreate(child, locationalParentGAI, node.GetGAI(), assets, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func collectAssetLikesAndDataToCreate(node AssetLike, locationalParentGAI, functionalParentGAI string, assets *[]AssetLikeWithParentReferences, data *[]Data, ts *time.Time, clientReference *string, visited map[string]bool) error {
	if visited[node.GetGAI()] {
		// Update existing relationships instead of skipping entirely
		for _, asset := range *assets {
			if asset.GetGAI() == node.GetGAI() {
				// Update locational and functional parent relationships
				if locationalParentGAI != "" && asset.GetLocationalParentGAI() == "" {
					asset.(*AssetLikeToCreate).locationalParentGAI = locationalParentGAI
				}
				if functionalParentGAI != "" && asset.GetFunctionalParentGAI() == "" {
					asset.(*AssetLikeToCreate).functionalParentGAI = functionalParentGAI
				}
				break
			}
		}
		return nil
	}
	visited[node.GetGAI()] = true

	assetID, err := node.GetAssetID()
	if err != nil {
		return fmt.Errorf("getting asset ID: %v", err)
	}

	if assetID == nil {
		a := &AssetLikeToCreate{
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
			err := collectAssetLikesAndDataToCreate(child, node.GetGAI(), functionalParentGAI, assets, data, ts, clientReference, visited)
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
			err := collectAssetLikesAndDataToCreate(child, locationalParentGAI, node.GetGAI(), assets, data, ts, clientReference, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AssetLikeToCreate implements AssetLikeWithParentReferences and holds asset information for creation.
type AssetLikeToCreate struct {
	node                AssetLike
	locationalParentGAI string
	functionalParentGAI string
}

func (a *AssetLikeToCreate) GetName() string {
	return a.node.GetName()
}

func (a *AssetLikeToCreate) GetDescription() string {
	return a.node.GetDescription()
}

func (a *AssetLikeToCreate) GetAssetType() string {
	return a.node.GetAssetType()
}

func (a *AssetLikeToCreate) GetGAI() string {
	return a.node.GetGAI()
}

func (a *AssetLikeToCreate) GetLocationalParentGAI() string {
	return a.locationalParentGAI
}

func (a *AssetLikeToCreate) GetFunctionalParentGAI() string {
	return a.functionalParentGAI
}

func (a *AssetLikeToCreate) GetSiteID() string {
	return a.node.GetSiteID()
}

func (a *AssetLikeToCreate) SetAssetID(assetID int32) error {
	return a.node.SetAssetID(assetID)
}
