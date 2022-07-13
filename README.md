# go-eliona #

go-eliona is a Go library for accessing resources within an [eliona](https://www.eliona.io/) environment.

The library contains handy functions to access database resources, Kafka topics and eliona API endpoints.
Besides, the library offers useful tool like logging.

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

## Usage ##
 
- [Logging](log) for logging purposes
- [Database](db) to access databases
- [Http](http) to handle web requests
- [Kafka](kafka) to handle kafka topics
- [Api](api) to access the [Eliona Api](https://github.com/eliona-smart-building-assistant/eliona-api) (generated)
- [Eliona](eliona) to handle specific use cases
  - [App](eliona/app) functions for apps and patches
  - [Asset](eliona/asset) asset and asset type management 
  - [Dashboard](eliona/dashboard) functions for dashboards