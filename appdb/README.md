# go-eliona App DB 
The go-eliona Apps package provides handy methods to handle app in an eliona environment. This package
uses a direct connection to the Eliona DB.

## Installation
To use the apps package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/appdb"
```

The package needs a database connection. To create and configure a database connection
have a look at [database](../db) package.

## Usage

After installation, you can initialize and patch apps in an eliona environment. 

```go
import "github.com/eliona-smart-building-assistant/go-eliona/appdb"
```

### Initialize your app

For the first start the app should be initialized. You can create your own schema and database tables to
persist configuration data for the app. Or you create example data that shows users how your app works.
To do this, you can use the `Init` function. This function is called only one time. After this the `Init` function skips all executions.

```go
appdb.Init(db.Pool(), common.AppName(),
    app.ExecSqlFile("database/init.sql"),
    app.ExecSqlFile("database/defaults.sql"))
```

You should call the `Init` at top in `main()` function.

### Patching your app

If you need to change data models, configuration tables or other things, you have to patch your app. That guarantees that installed apps can always be updated even though they have already been initialized. To do this, you can use the `Patch` function. This function is called once for each patch. After this the `Patch` function skips all executions for this patch.

```go
appdb.Patch(db.Pool(), common.AppName(), "010100",
    app.ExecSqlFile("database/patches/010100.sql"))
```

You should call the `Patch` at top in `main()` after the `Init` function.