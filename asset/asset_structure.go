package asset

import (
	"fmt"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
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
		traverseCreated, err := traverseFunctionalTree(fc, projectId, rootAssetID, rootAssetID)
		if err != nil {
			return createdCnt, fmt.Errorf("functional tree traversal: %v", err)
		}
		createdCnt += traverseCreated
	}

	for _, lc := range root.GetLocationalChildren() {
		if lc == nil {
			continue
		}
		traverseCreated, err := traverseLocationalTree(lc, projectId, rootAssetID, rootAssetID)
		if err != nil {
			return createdCnt, fmt.Errorf("locational tree traversal: %v", err)
		}
		createdCnt += traverseCreated
	}
	return createdCnt, nil
}

func traverseLocationalTree(node LocationalNode, projectId string, locationalParentAssetId, functionalParentAssetId *int32) (createdCnt int, err error) {
	currentAssetId, created, err := createAsset(node, projectId, locationalParentAssetId, functionalParentAssetId)
	if err != nil {
		return createdCnt, err
	}
	if created {
		createdCnt++
	}

	for _, child := range node.GetLocationalChildren() {
		if child == nil {
			continue
		}
		traverseCreated, err := traverseLocationalTree(child, projectId, currentAssetId, functionalParentAssetId)
		if err != nil {
			return createdCnt, err
		}
		createdCnt += traverseCreated
	}
	return createdCnt, nil
}

func traverseFunctionalTree(node FunctionalNode, projectId string, locationalParentAssetId, functionalParentAssetId *int32) (createdCnt int, err error) {
	currentAssetId, created, err := createAsset(node, projectId, locationalParentAssetId, functionalParentAssetId)
	if err != nil {
		return createdCnt, err
	}
	if created {
		createdCnt++
	}

	for _, child := range node.GetFunctionalChildren() {
		if child == nil {
			continue
		}
		traverseCreated, err := traverseFunctionalTree(child, projectId, locationalParentAssetId, currentAssetId)
		if err != nil {
			return createdCnt, err
		}
		createdCnt += traverseCreated
	}
	return createdCnt, nil
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
		Id:                      *api.NewNullableInt32(originalAssetID),
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
