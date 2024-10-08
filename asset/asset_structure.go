package asset

import (
	"fmt"
	"time"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v2/tools"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
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

func CreateAssets(root Root, projectId string) (createdCnt int, err error) {
	return createAssetsAndUpsertData(root, projectId, false, nil, nil)
}

func CreateAssetsAndUpsertData(root Root, projectId string, ts *time.Time, clientReference *string) (createdCnt int, err error) {
	return createAssetsAndUpsertData(root, projectId, true, ts, clientReference)
}

func createAssetsAndUpsertData(root Root, projectId string, upsertingData bool, ts *time.Time, clientReference *string) (createdCnt int, err error) {
	rootAssetID, created, err := createRoot(root, projectId)
	if err != nil {
		return createdCnt, fmt.Errorf("upserting root asset: %v", err)
	}
	if created {
		createdCnt++
	}
	for _, fc := range root.GetFunctionalChildren() {
		if fc == nil {
			continue
		}
		traverseCreated, err := traverseFunctionalTree(fc, projectId, rootAssetID, rootAssetID, upsertingData, ts, clientReference)
		if err != nil {
			return createdCnt, fmt.Errorf("functional tree traversal: %v", err)
		}
		createdCnt += traverseCreated
	}

	for _, lc := range root.GetLocationalChildren() {
		if lc == nil {
			continue
		}
		traverseCreated, err := traverseLocationalTree(lc, projectId, rootAssetID, rootAssetID, upsertingData, ts, clientReference)
		if err != nil {
			return createdCnt, fmt.Errorf("locational tree traversal: %v", err)
		}
		createdCnt += traverseCreated
	}
	return createdCnt, nil
}

func traverseLocationalTree(
	node LocationalNode,
	projectId string,
	locationalParentAssetId,
	functionalParentAssetId *int32,
	upsertingData bool,
	ts *time.Time,
	clientReference *string,
) (createdCnt int, err error) {

	currentAssetId, created, err := createAsset(node, projectId, locationalParentAssetId, functionalParentAssetId)
	if err != nil {
		return createdCnt, err
	}
	if !created {
		if err := setAssetParent(currentAssetId, node, projectId, locationalParentAssetId, nil); err != nil {
			return createdCnt, fmt.Errorf("updating asset parent: %v", err)
		}
	}
	if created {
		createdCnt++
	}

	if currentAssetId != nil && upsertingData {
		err = upsertNodeDataIfAssetExists(node, *currentAssetId, ts, clientReference)
		if err != nil {
			return createdCnt, err
		}
	}

	for _, child := range node.GetLocationalChildren() {
		if child == nil {
			continue
		}
		traverseCreated, err := traverseLocationalTree(child, projectId, currentAssetId, functionalParentAssetId, upsertingData, ts, clientReference)
		if err != nil {
			return createdCnt, err
		}
		createdCnt += traverseCreated
	}
	return createdCnt, nil
}

func traverseFunctionalTree(
	node FunctionalNode,
	projectId string,
	locationalParentAssetId,
	functionalParentAssetId *int32,
	upsertingData bool,
	ts *time.Time,
	clientReference *string,
) (createdCnt int, err error) {

	currentAssetId, created, err := createAsset(node, projectId, locationalParentAssetId, functionalParentAssetId)
	if err != nil {
		return createdCnt, err
	}
	if !created {
		if err := setAssetParent(currentAssetId, node, projectId, nil, functionalParentAssetId); err != nil {
			return createdCnt, fmt.Errorf("updating asset parent: %v", err)
		}
	}
	if created {
		createdCnt++
	}

	if currentAssetId != nil && upsertingData {
		err = upsertNodeDataIfAssetExists(node, *currentAssetId, ts, clientReference)
		if err != nil {
			return createdCnt, err
		}
	}

	for _, child := range node.GetFunctionalChildren() {
		if child == nil {
			continue
		}
		traverseCreated, err := traverseFunctionalTree(child, projectId, locationalParentAssetId, currentAssetId, upsertingData, ts, clientReference)
		if err != nil {
			return createdCnt, err
		}
		createdCnt += traverseCreated
	}
	return createdCnt, nil
}

func upsertNodeDataIfAssetExists(node Asset, assetId int32, ts *time.Time, clientReference *string) error {
	cr := ""
	if clientReference != nil {
		cr = *clientReference
	}
	return UpsertAssetDataIfAssetExists(Data{
		AssetId:         assetId,
		Timestamp:       *api.NewNullableTime(ts),
		ClientReference: cr,
		Data:            node,
	})
}

func createRoot(ast Asset, projectId string) (assetId *int32, created bool, err error) {
	return createAsset(ast, projectId, nil, nil)
}

func createAsset(ast Asset, projectId string, locationalParentAssetId *int32, functionalParentAssetId *int32) (assetId *int32, created bool, err error) {
	originalAssetID, err := ast.GetAssetID(projectId)
	if err != nil {
		return nil, created, fmt.Errorf("getting asset id: %v", err)
	}
	if originalAssetID != nil {
		return originalAssetID, created, nil
	}
	a := api.Asset{
		ProjectId:               projectId,
		GlobalAssetIdentifier:   ast.GetGAI(),
		Name:                    *api.NewNullableString(common.Ptr(ast.GetName())),
		AssetType:               ast.GetAssetType(),
		Description:             *api.NewNullableString(common.Ptr(ast.GetDescription())),
		ParentFunctionalAssetId: *api.NewNullableInt32(functionalParentAssetId),
		ParentLocationalAssetId: *api.NewNullableInt32(locationalParentAssetId),
	}
	assetID, err := UpsertAsset(a)
	if err != nil {
		return nil, created, fmt.Errorf("upserting asset %+v into Eliona: %v", a, err)
	}
	if assetID == nil {
		return nil, created, fmt.Errorf("cannot create asset %s", ast.GetName())
	}
	created = true

	if err := ast.SetAssetID(*assetID, projectId); err != nil {
		return nil, created, fmt.Errorf("setting asset id: %v", err)
	}
	return assetID, created, nil
}

func setAssetParent(assetId *int32, ast Asset, projectId string, locationalParentAssetId *int32, functionalParentAssetId *int32) error {
	asset := api.Asset{
		Id:                    *api.NewNullableInt32(assetId),
		ProjectId:             projectId,
		GlobalAssetIdentifier: ast.GetGAI(),
		AssetType:             ast.GetAssetType(),
	}
	fetchedAsset, _, err := client.NewClient().AssetsAPI.
		GetAssetById(client.AuthenticationContext(), asset.GetId()).
		Execute()
	if err != nil {
		tools.LogError(fmt.Errorf("getting asset %v: %w", asset.GetId(), err))
	}
	if locationalParentAssetId != nil && !fetchedAsset.HasParentLocationalAssetId() {
		asset.ParentLocationalAssetId = *api.NewNullableInt32(locationalParentAssetId)
	}
	if functionalParentAssetId != nil && !fetchedAsset.HasParentFunctionalAssetId() {
		asset.ParentFunctionalAssetId = *api.NewNullableInt32(functionalParentAssetId)
	}

	// TODO: This is an ugly workaround to not overwrite name changes in Eliona.
	// We have to get rid of this once APIv2 starts skipping empty fields.
	if !asset.Name.IsSet() {
		asset.Name = fetchedAsset.Name
	}
	if !asset.Description.IsSet() {
		asset.Description = fetchedAsset.Description
	}
	//

	if _, err := UpsertAsset(asset); err != nil {
		return fmt.Errorf("upserting asset %+v into Eliona: %v", asset, err)
	}
	return nil
}
