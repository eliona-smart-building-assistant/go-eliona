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

package appdb

import (
	"context"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
)

// The Init function must be used to run all the elements required for the app initialization process.
// This function guarantees that everything will only run once when the app is first launched.
// Furthermore, this function guarantees that either all database changes or no changes are committed using
// transactions. For this you must use the connection that is passed to the function parameter.
func Init(connection db.Connection, appName string, initFunctions ...func(connection db.Connection) error) {
	if appInitialized(connection, appName) {
		log.Debug("Apps", "App %s is already initialized. Skip init.", appName)
		return
	}

	transaction, err := db.Begin(connection)
	if err != nil {
		log.Fatal("Apps", "Cannot start transaction to init app %s.", appName)
	}

	for i, initFunction := range initFunctions {
		err := initFunction(transaction)
		if err != nil {
			log.Fatal("Apps", "Cannot execute function %d to init app %s.", i, appName)
		}
	}

	err = registerApp(transaction, appName)
	if err != nil {
		log.Fatal("Apps", "Cannot register app %s as initialized.", appName)
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		log.Fatal("Apps", "Cannot commit init for app %s.", appName)
	}
}

// appInitialized checks if the app is already initialized.
func appInitialized(connection db.Connection, appName string) bool {
	count, _ := db.QuerySingleRow[int](connection, "select count(*) from public.eliona_app where app_name = $1 and initialised", appName)
	return 0 < count
}

// registerApp marks that the app is now initialized and installed.
func registerApp(connection db.Connection, appName string) error {
	err := db.Exec(connection, "insert into public.eliona_app (app_name, category, active, initialised) values"+
		" ($1, 'app', true, true) "+
		" on conflict (app_name) do update set initialised = true", appName)
	if err != nil {
		log.Error("Apps", "Cannot register app %s.", appName)
		return err
	}
	return nil
}

// The Patch function must be used to run all the elements required for the patch process.
// This function guarantees that everything will only run once when the patch is applied.
// Furthermore, this function guarantees that either all database changes or no changes are committed using
// transactions. For this you must use the connection that is passed to the function parameter.
func Patch(connection db.Connection, appName string, patchName string, patchFunctions ...func(connection db.Connection) error) {
	if patchRegistered(connection, appName, patchName) {
		log.Debug("Apps", "App %s patch %s is already installed. Skip patching.", appName, patchName)
		return
	}

	transaction, err := db.Begin(connection)
	if err != nil {
		log.Fatal("Apps", "Cannot start transaction to patch %s app %s.", patchName, appName)
	}

	for i, patchFunction := range patchFunctions {
		err := patchFunction(transaction)
		if err != nil {
			log.Fatal("Apps", "Cannot execute function %d to patch %s app %s.", i, patchName, appName)
		}
	}

	err = registerPatch(transaction, appName, patchName)
	if err != nil {
		log.Fatal("Apps", "Cannot register patch %s for app %s.", patchName, appName)
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		log.Fatal("Apps", "Cannot commit patch %s for app %s.", patchName, appName)
	}
}

// patchRegistered checks if the patch is already applied.
func patchRegistered(connection db.Connection, appName string, patchName string) bool {
	count, _ := db.QuerySingleRow[int](connection, "select count(*) from versioning.patches where app_name = $1 and patch_name = $2", appName, patchName)
	return 0 < count
}

// registerPatch marks that the patch is now applied.
func registerPatch(connection db.Connection, appName string, patchName string) error {
	err := db.Exec(connection, "insert into versioning.patches "+
		"(patch_name, app_name, applied_tsz, applied_by) values "+
		"($1, $2, now(), current_user)", patchName, appName)
	if err != nil {
		log.Error("Apps", "Cannot register patch %s for app %s.", patchName, appName)
		return err
	}
	return nil
}
