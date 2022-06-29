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
- [API](api) to access API (generated, see [instruction for generation](api/README.md))
