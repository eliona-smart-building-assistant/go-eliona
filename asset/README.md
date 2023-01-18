# go-eliona Asset DB
The go-eliona Assets package provides functions and data structures to handle assets and data. This package uses the [Eliona API](https://github.com/eliona-smart-building-assistant/eliona-api) to access Eliona.

## Installation
To use the assets package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/asset"
```
The package needs an Api server address. To create and configure have a look at [api](../api) package.

## Usage

After installation, you can use the assets package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/asset"
```

### Configuring asset types and attributes

You can create new assets types and attributes or change existing ones. For example, if you want to create a weather location asset type that holds temperature data, you have to create the following.

```go
_ = asset.UpsertAssetType(api.AssetType{Name: "weather_location", Custom: true, Vendor: "ITEC AG", Translation: api.Translation{De: "Wetterstation", En: "Weather location"}})
```

### Write asset data

For example, you can insert or update asset data for temperatures of type `Temperature`. To do this, you can use the defined `Data` data structure with data field.

```go
type Temperature struct {
    Value int
    Unit  string
}
```

The following code uses the `UpsertData()` function and inserts a data for
asset with id `2` and  the `info` subtype. If already exists, the
data and timestamp are updated. The data has 'now' as timestamp and a temperature with `23` as value and `Celsius` as unit. This would be written as `{"Unit": "Celsius", "Value": 30}` to data payload.

```go
_ = asset.UpsertData(api.Data{2, api.INFO, time.Time{}, common.StructToMap(Temperature{35, "Celsius"})})
```
