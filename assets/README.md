# go-eliona Assets
The go-eliona Assets package provides functions and data structures to handle assets and heaps.

## Installation
To use the assets package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/assets"
```

The package needs a database connection. To create and configure a database connection
have a look at [database](../db) package.    

## Usage

After installation, you can use the assets package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/assets"
import "github.com/eliona-smart-building-assistant/go-eliona/db"
```

### Configuring asset types and attributes

You can create new assets types and attributes or change existing ones. For example, if you want to create a weather location asset type that holds temperature data, you have to create the following.

```go
_ = assets.UpsertAssetType(db.Pool(), assets.AssetType{Id: "weather_location", Custom: true, Vendor: "ITEC AG", Translation: Translation{German: "Wetterstation", English: "Weather location"}})
_ = assets.UpsertAssetTypeAttribute(db.Pool(), assets.AssetTypeAttribute{AssetTypeId: "weather_location", AttributeType: "temperature", Id: "temperature", Subtype: assets.InputSubtype, Enable: true, Translation: Translation{German: "Temperatur", English: "Temperature"}})
```

### Write a heap

For example, you can insert or update a heap for temperatures. To do this, you can use the defined `Heap` data structure which
has a generic type for the data. If the data corresponds to `Temperature` structure, so you have to use `Heap[Temperature]`.

```go
type Temperature struct {
    Value int
    Unit  string
}
```

The following code uses the `UpsertHeap()` function and inserts a heap for
asset with id `2` and  the `info` subtype. If already exists, the
data and timestamp are updated. The heap has 'now' as timestamp and a temperature with `23` as value
and `Celsius` as unit. This would be written as `{"Unit": "Celsius", "Value": 30}` to heap data.

```go
_ = assets.UpsertHeap(db.Pool(), assets.Heap[Temperature]{2, InfoSubtype, time.Time{}, Temperature{35, "Celsius"}})
```

### Listen for changed heaps

The eliona environment informs, if heaps are inserted and updated. You can use this to get e.g. set points for
assets. Usually set points are defines as `OutputSubtype` of heaps.

For example, you have to get a set point for temperatures. Then you can listen for changes on heap data and use this
to modify the setting of assets. To do this, use the `ListenHeap` function and a go channel with the defined `Heap`
data structure including the `Temperature` structure (`Heap[Temperature]`).

```go
heaps := make(chan assets.Heap[Temperature])
go assets.ListenHeap(db.Pool(), heaps)
for heap := range heaps {
    if heap.Subtype == assets.OutputSubtype {
        log.Debug("Assets", "Do something for asset with id %d and value %d unit %s.", heap.AssetId, heap.GetData().Value, heap.GetData().Unit)
    }
}
```
