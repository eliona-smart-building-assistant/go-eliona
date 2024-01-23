# Go Eliona

Go Eliona is a part of the Eliona App SDK. It provides a Go library for accessing resources within an Eliona environment.

## Installation ##

To get go-eliona you can use command line:

```bash
go get github.com/eliona-smart-building-assistant/go-eliona
```

or you define import in go files:

```go
import "github.com/eliona-smart-building-assistant/go-eliona"
```

and run `go get` without parameters.

## Configuration

The `API_ENDPOINT` variable configures the endpoint to access the [Eliona API](https://github.com/eliona-smart-building-assistant/eliona-api). If the app runs as a Docker container inside an Eliona environment, the environment must provide this variable. If you run the app standalone you must set this variable. Otherwise, the app can't be initialized and started. The default endpoint is `http://api-v2:3000/v2`

```bash
export API_ENDPOINT=http://api-v2:3000/v2
```

The `API_TOKEN` defines the secret to authenticate the app and access the API.

## Usage ##
 
- [App](app) functions for apps and patches
- [Asset](asset) asset and asset type management 
- [Dashboard](dashboard) functions for dashboards
- [Frontend](frontend) functions for frontend
