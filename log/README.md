# go-eliona Logging 
The go-eliona Logging package provides handy methods for standardized logging.

## Installation
To use the log package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/log"
```

Optionally you can define an environment variable named LOG_LEVEL which sets the maximum level to be logged.
Not defined the log packages takes the info log level as default.

```bash
export LOG_LEVEL=debug # This is optionally, default is info log level
```

## Usage

After installation, you can use the logging like this:

```go
import "github.com/eliona-smart-building-assistant/go-eliona/log"

log.Info("start", "This could interest %d of %d people", 3, 4)
log.Debug("start", Usually nobody cares %s", "this")
log.Fatal("validate", "The house is on fire!") // exits with code 1
```

This produces the following output:

```
INFO    2022-05-03 15:34:41.770515      start   This could interest 3 of 4 people
DEBUG   2022-05-03 15:34:41.749597      start   Usually nobody cares this
FATAL   2022-05-03 15:34:41.770515      validate   The house is on fire!
```
