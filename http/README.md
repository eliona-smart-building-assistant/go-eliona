# go-eliona Http 
The go-eliona Http package provides handy methods to use api end points using http protocol.

## Installation
To use the log package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/http"
```

## Usage

After installation, you can read payload from http end points.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/http"
```

For example, you can make a request and read current weather conditions from this endpoint.  

```go
payload, err := http.Read(http.NewRequest("https://weatherdbi.herokuapp.com/data/weather/winterthur"), 10, true)
var result map[string]interface{}
_ = json.Unmarshal(payload, &result)
fmt.Printf(result["currentConditions"].(map[string]interface{})["comment"].(string))
```