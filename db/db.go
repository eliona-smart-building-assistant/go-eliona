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

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

// ConnectionString returns the connection string defined in the environment variable CONNECTION_STRING.
// If not defined, the method returns postgres://user:secret@database/iot as default.
func ConnectionString() string {
	return common.Getenv("CONNECTION_STRING", "postgres://user:secret@database/iot")
}

// The Connection interface allows mocking database connection for testing
type Connection interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

// ConnectionConfig returns the connection config defined by CONNECTION_STRING and APPNAME environment variables.
// This is sufficient for the most common cases using connections within eliona apps.
func ConnectionConfig() *pgx.ConnConfig {
	config, err := pgx.ParseConfig(ConnectionString())
	if err != nil {
		log.Fatal("Database", "Unable to parse database URL: %v", err)
	}
	// add Application name tho
	config.RuntimeParams["application_name"] = common.AppName()
	return config
}

func PoolConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(ConnectionString())
	if err != nil {
		log.Fatal("Database", "Unable to parse database URL: %v", err)
	}
	return config
}

func ExecFile(connection Connection, path string) error {
	sql, err := ioutil.ReadFile(filepath.Join(path))
	if err != nil {
		log.Error("Database", "Unable to read sql file %s: %v", path, err)
		return err
	}
	_, err = connection.Exec(context.Background(), string(sql))
	if err != nil {
		log.Error("Database", "Error during execute sql file %s: %v", path, err)
		return err
	}
	return nil
}

// NewConnection returns a new connection defined by CONNECTION_STRING environment variable.
// The new connection can be hold by apps and must be closed independently.
func NewConnection() *pgx.Conn {
	connection, err := pgx.ConnectConfig(context.Background(), ConnectionConfig())
	if err != nil {
		log.Fatal("Database", "Unable to create connection to database: %v", err)
	}
	log.Debug("Database", "Connection created")
	return connection
}

func NewPool() *pgxpool.Pool {
	pool, err := pgxpool.ConnectConfig(context.Background(), PoolConfig())
	if err != nil {
		log.Fatal("Database", "Unable to create pool for database: %v", err)
	}
	log.Debug("Database", "Pool created")
	return pool
}

// current holds a single connection
var poolMutex sync.Mutex
var pool *pgxpool.Pool

// Pool returns the default pool hold by this package. The pool is created if this function is called first time.
// Afterwards this function returns always the same pool. Don't forget to defer the pool with ClosePool function.
func Pool() *pgxpool.Pool {
	if pool == nil {
		poolMutex.Lock()
		if pool == nil {
			pool = NewPool()
		}
		poolMutex.Unlock()
	}
	return pool
}

// ClosePool closes the default pool hold by this package.
func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}

// Listen waits for notifications on database channel and writes the payload to the go channel.
// The type of the go channel have to correspond to the payload JSON structure
func Listen[T any](connection *pgx.Conn, channel string, payloads chan T) {
	contextWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := connection.Exec(contextWithCancel, "LISTEN "+channel)
	if err != nil {
		log.Error("Database", "Error listening on channel '%s': %v", channel, err)
	}
	for {
		notification, _ := waitForNotification(contextWithCancel, connection)
		if notification != nil {
			var payload T
			err := json.Unmarshal([]byte(notification.Payload), &payload)
			if err != nil {
				log.Error("Database", "Unmarshal error during listening: %v", err)
			}
			payloads <- payload
		}
	}
}

// waitForNotification waits for channel notification of the given connection
func waitForNotification(origCtx context.Context, connection *pgx.Conn) (*pgconn.Notification, error) {
	ctx, cancel := context.WithTimeout(origCtx, 5*time.Second)
	defer cancel()
	notification, err := connection.WaitForNotification(ctx)
	if err == nil {
		return notification, nil
	} else if pgconn.Timeout(err) {
		ctx, cancel = context.WithTimeout(origCtx, 1*time.Second)
		defer cancel()
		err = connection.Ping(ctx)
	}
	if err != nil {
		log.Error("Database", "Error waiting for notification: %v", err)
	}
	return nil, err
}

// Exec inserts a row using the given sql with arguments
func Exec(connection Connection, sql string, args ...interface{}) error {
	_, err := connection.Exec(context.Background(), sql, args...)
	if err != nil {
		log.Error("Database", "Error in statement '%s': %v", sql, err)
	}
	return err
}

// EmptyStringIsNull returns a null string if the string is empty (string == "")
func EmptyStringIsNull(string string) sql.NullString {
	if len(string) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: string, Valid: true}
}

func EmptyFloatIsNull(float float64) sql.NullFloat64 {
	if float == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: float, Valid: true}
}

// Begin returns a new transaction
func Begin(connection Connection) (pgx.Tx, error) {
	transaction, err := connection.Begin(context.Background())
	if err != nil {
		log.Error("Database", "Error starting transaction: %v", err)
		return transaction, err
	}
	return transaction, nil
}

// Query gets values read from database into a channel. The value type of channel must match
// the fields defined in the query. The type can be a single value (e.g. string) if the query
// returns only a single field. Otherwise, the type have to be a struct with the identical number
// of elements and corresponding types like the query statement
func Query[T any](connection Connection, sql string, results chan T, args ...interface{}) error {
	defer close(results)
	rows, err := connection.Query(context.Background(), sql, args...)
	if err != nil {
		log.Error("Database", "Error in query statement '%s': %v", sql, err)
		return err
	} else {
		defer rows.Close()
		for rows.Next() {
			var result T
			err := rows.Scan(interfaces(&result)...)
			if err != nil {
				log.Error("Database", "Error scanning result '%s': %v", sql, err)
				return err
			}
			results <- result
		}
	}
	return nil
}

// QuerySingleRow returns the value if only a single row is queried
func QuerySingleRow[T any](connection Connection, sql string, args ...interface{}) (T, error) {
	result := make(chan T)
	err := make(chan error)
	defer close(err)
	go func() {
		err <- Query(connection, sql, result, args...)
	}()
	return <-result, <-err
}

// interfaces creates interface for the given holder, or if the holder a structure this
// function returns a list of interfaces for all structure members.
func interfaces(holder interface{}) []interface{} {
	value := reflect.ValueOf(holder).Elem()
	if value.Kind() == reflect.Struct {
		values := make([]interface{}, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			values[i] = value.Field(i).Addr().Interface()
		}
		return values
	} else {
		values := make([]interface{}, 1)
		values[0] = value.Addr().Interface()
		return values
	}
}
