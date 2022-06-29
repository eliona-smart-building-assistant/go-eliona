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

package assetdb

import (
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"math"
)

// Translation defines a translation used inside assetdb
type Translation struct {
	German  string `json:"de,omitempty"`
	English string `json:"en,omitempty"`
}

type Pipeline struct {
	Mode   PipelineMode `json:"mode,omitempty"`
	Raster string       `json:"raster,omitempty"`
}

// AssetType defines an asset type
type AssetType struct {
	Name             string               `json:"name,omitempty"`
	Custom           bool                 `json:"custom,omitempty"`
	Vendor           string               `json:"vendor,omitempty"`
	Translation      *Translation         `json:"translation,omitempty"`
	DocumentationUrl string               `json:"urldoc,omitempty"`
	Icon             string               `json:"icon,omitempty"`
	Attributes       []AssetTypeAttribute `json:"attributes,omitempty"`
}

// UpsertAssetType insert or, when already exist, updates an asset type
func UpsertAssetType(connection db.Connection, assetType AssetType) error {
	err := db.Exec(connection,
		"insert into public.asset_type ("+
			"asset_type,"+
			"custom,"+
			"vendor,"+
			"translation,"+
			"urldoc,"+
			"icon"+
			") values ($1, $2, $3, $4, $5, $6) "+
			"on conflict (asset_type) "+
			"do update set custom = excluded.custom, vendor = excluded.vendor, translation = excluded.translation, urldoc = excluded.urldoc, icon = excluded.icon",
		assetType.Name,
		assetType.Custom,
		db.EmptyStringIsNull(&assetType.Vendor),
		db.EmptyJsonIsNull(assetType.Translation),
		db.EmptyStringIsNull(&assetType.DocumentationUrl),
		db.EmptyStringIsNull(&assetType.Icon),
	)
	if err != nil {
		return err
	}
	if assetType.Attributes != nil {
		for _, attribute := range assetType.Attributes {
			attribute.AssetTypeId = assetType.Name
			err = UpsertAssetTypeAttribute(connection, attribute)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AssetTypeAttribute defines an attribute for asset type
type AssetTypeAttribute struct {
	AssetTypeId   string       `json:"assetType,omitempty"`
	AttributeType string       `json:"type,omitempty"`
	Name          string       `json:"name,omitempty"`
	Subtype       Subtype      `json:"subtype,omitempty"`
	Enable        bool         `json:"enable,omitempty"`
	Translation   *Translation `json:"translation,omitempty"`
	Unit          string       `json:"unit,omitempty"`
	Pipeline      Pipeline     `json:"pipeline,omitempty"`
	Precision     *int16       `json:"precision,omitempty"`
}

// UpsertAssetTypeAttribute insert or, when already exist, updates an attribute for asset type
func UpsertAssetTypeAttribute(connection db.Connection, attribute AssetTypeAttribute) error {
	return db.Exec(connection,
		"insert into public.attribute_schema ("+
			"asset_type,"+
			"attribute_type,"+
			"attribute,"+
			"subtype,"+
			"enable,"+
			"translation,"+
			"unit,"+
			"pipeline_mode,"+
			"pipeline_raster,"+
			"precision"+
			") values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) "+
			"on conflict (asset_type, subtype, attribute) "+
			"do update set attribute_type = excluded.attribute_type, enable = excluded.enable, translation = excluded.translation, unit = excluded.unit, "+
			"pipeline_mode = excluded.pipeline_mode, pipeline_raster = excluded.pipeline_raster, precision = excluded.precision",
		attribute.AssetTypeId,
		attribute.AttributeType,
		attribute.Name,
		string(attribute.Subtype),
		attribute.Enable,
		db.EmptyJsonIsNull(attribute.Translation),
		db.EmptyStringIsNull(&attribute.Unit),
		db.EmptyStringIsNull(&attribute.Pipeline.Mode),
		db.EmptyStringIsNull(&attribute.Pipeline.Raster),
		db.SmallIntIsNull(attribute.Precision, math.MinInt16),
	)
}

// Asset defines an asset
type Asset struct {
	ProjectId             string   `json:"projectId"`
	GlobalAssetIdentifier string   `json:"gai"`
	Name                  string   `json:"name"`
	AssetTypeId           string   `json:"asset_type"`
	Latitude              float64  `json:"lat"`
	Longitude             float64  `json:"lon"`
	Description           string   `json:"description"`
	Tags                  []string `json:"tags"`
}

// ExistAsset returns true, if the given asset id exists in eliona
func ExistAsset(connection db.Connection, assetId int) (bool, error) {
	count, err := db.QuerySingleRow[int](connection,
		"select count(*) from public.asset where asset_id = $1", assetId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpsertAsset insert or updates an assetdb and returns the id
func UpsertAsset(connection db.Connection, asset Asset) (int, error) {
	assetId, err := db.QuerySingleRow[int](connection,
		"with asset_id as ("+
			"insert into public.asset ("+
			"proj_id,"+
			"gai,"+
			"name,"+
			"asset_type,"+
			"lat,"+
			"lon,"+
			"description,"+
			"tags"+
			") values ($1, $2, $3, $4, $5, $6, $7, $8) "+
			"on conflict (proj_id, gai) "+
			"do update set name = excluded.name, asset_type = excluded.asset_type, lat = excluded.lat, lon = excluded.lon, description = excluded.description, tags = excluded.tags "+
			"returning asset_id "+
			") select asset_id from asset_id",
		asset.ProjectId,
		asset.GlobalAssetIdentifier,
		asset.Name,
		asset.AssetTypeId,
		db.EmptyFloatIsNull(&asset.Latitude),
		db.EmptyFloatIsNull(&asset.Longitude),
		db.EmptyStringIsNull(&asset.Description),
		asset.Tags,
	)
	return assetId, err
}

// PipelineMode defines the type of pipeline
type PipelineMode string

const (
	// NoPipelineMode defines, that no pipeline mode ist used
	NoPipelineMode PipelineMode = ""

	// SumPipelineMode is the sum pipeline mode
	SumPipelineMode = "sum"

	// AveragePipelineMode is the a pipeline mode
	AveragePipelineMode = "avg"
)

// Subtype defines the subtype of heaps which is e.g. input or info
type Subtype string

const (
	// InputSubtype is the subtype info
	InputSubtype Subtype = ""

	// OutputSubtype is the subtype output
	OutputSubtype = "sp"

	// InfoSubtype is the subtype info
	InfoSubtype = "info"

	// StatusSubtype is the subtype status
	StatusSubtype = "status"
)