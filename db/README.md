# go-eliona Database 
The go-eliona Database package provides handy methods to read and write data from and to an
eliona database. The database uses [PostgreSQL](https://www.postgresql.org/) as management system.

## Installation
To use the log package you must import the package.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/db"
```

Optionally you can define an environment variable named CONNECTION_STRING which defines the database
that should be connected to. In eliona environment the CONNECTION_STRING is set by default to the internal database `database/iot`
 and the user `app`.

```bash
export CONNECTION_STRING=postgres://user:pass@localhost/iot # This is optionally, default is the internal eliona database
```

## Usage

After installation, you can access the eliona database.

```go
import "github.com/eliona-smart-building-assistant/go-eliona/db"
```
SQL statements have to apply the guidelines of [PostgreSQL](https://www.postgresql.org/). Placeholders are defined as `$1`, `$2` and so on.

### Reading data from database

For example, you can read temperature objects based on an SQL statement.
You have to open a new connection and read the result through a channel.

```go
type Temperature struct {
    Value int
    Unit  string
}
```

Note, that the fields of `Temperature` correspond to the selected values in the SQL statement. Here `Value` belongs to
the constant numeric value `23` and `Unit` to constant string `Celsius`.

```go
connection := db.NewConnection()
defer connection.Close(context.Background())
temperatures := make(chan Temperature)
go db.Query(connection, "select 23, 'Celsius'", temperatures)
for temperature := range temperatures {
    log.Debug("Temperature", "Temperature is: %d %s", temperature.Value, temperature.Unit)
}
```

To select single values only, you can use simple typed channels instead of structures.

```go
connection := db.NewConnection()
defer connection.Close(context.Background())
values := make(chan int)
go db.Query(connection, "select 23", values)
for value := range values {
	log.Debug("Temperature", "Value is: %d", value)
}
```

### Modifying data in database

For example, you can write a temperature in a database table named `temperatures`.

```go
connection := db.NewConnection()
defer connection.Close(context.Background())
_ = db.Exec(connection, "insert into temperatures (value, unit) values ($1, $2)",
	23, "Celsius")
```

In the same way with `db.Exec()` you can update, delete and modifying data as you want.

### Listen on database channels

You can listen for notifications on database channels. For example, a notification is triggered
for changes on a database table with temperatures, so the database channel is named `temperatures`.
The payload of the notification have to be a json string, e.g. `{"Value": 30, "Unit": "Celsius"}`. 
In sum the notification is triggered as `pg_notify('temperatures', '{"Value": 30, "Unit": "Celsius"}')`.

You can use a go channel to read this into a corresponding structure. For example notification we can use
the `Temperature` structure above. The notification payload will be mapped to this structure.  

```go
connection := db.NewConnection()
temperatures := make(chan Temperature)
go db.Listen(connection, "temperatures", temperatures)
for temperature := range temperatures {
    log.Debug("Temperature", "Temperature is: %d %s", temperature.Value, temperature.Unit)
}
```