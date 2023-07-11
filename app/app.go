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

package app

import (
	"context"
	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v2/tools"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// AppName returns the name of the app uses the library. The app name is defined in the
// environment variable APPNAME. If not defined, the AppName returns nil.
func AppName() string {
	return common.Getenv("APPNAME", "")
}

// The ExecSqlFile returns a function which executes the given sql file. This method can be used
// as parameter for the Init and Patch function.
func ExecSqlFile(path string) func(connection db.Connection) error {
	return func(connection db.Connection) error {
		return db.ExecFile(connection, path)
	}
}

// The Init function must be used to run all the elements required for the app initialization process.
// This function guarantees that everything will only run once when the app is first launched.
// Furthermore, this function guarantees that either all database changes or no changes are committed using
// transactions. For this you must use the connection that is passed to the function parameter.
func Init(connection db.Connection, appName string, initFunctions ...func(connection db.Connection) error) {
	if appRegistered(appName) {
		log.Debug("Apps", "App %s is already initialized. Skip init.", appName)
		return
	}

	transaction, err := db.Begin(connection)
	if err != nil {
		log.Fatal("Apps", "Cannot start transaction to init app %s: %v", appName, err)
	}

	for i, initFunction := range initFunctions {
		err := initFunction(transaction)
		if err != nil {
			log.Fatal("Apps", "Cannot execute function %d to init app %s: %v", i, appName, err)
		}
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		log.Fatal("Apps", "Cannot commit init for app %s: %v", appName, err)
	}

	err = registerApp(appName)
	if err != nil {
		log.Fatal("Apps", "Cannot register app %s as initialized: %v", appName, err)
	}

}

// appRegistered checks if the app is already initialized.
func appRegistered(appName string) bool {
	app, _, err := client.NewClient().AppsApi.
		GetAppByName(client.AuthenticationContext(), appName).
		Execute()
	tools.LogError(err)
	if err != nil || !app.Registered.IsSet() {
		return false
	}
	return *app.Registered.Get()
}

// registerApp marks that the app is now initialized and installed.
func registerApp(appName string) error {
	_, err := client.NewClient().AppsApi.
		PatchAppByName(client.AuthenticationContext(), appName).
		Registered(true).
		Execute()
	tools.LogError(err)
	return err
}

// The Patch function must be used to run all the elements required for the patch process.
// This function guarantees that everything will only run once when the patch is applied.
// Furthermore, this function guarantees that either all database changes or no changes are committed using
// transactions. For this you must use the connection that is passed to the function parameter.
func Patch(connection db.Connection, appName string, patchName string, patchFunctions ...func(connection db.Connection) error) {
	if patchApplied(appName, patchName) {
		log.Debug("Apps", "App %s patch %s is already installed. Skip patching.", appName, patchName)
		return
	}

	transaction, err := db.Begin(connection)
	if err != nil {
		log.Fatal("Apps", "Cannot start transaction to patch %s app %s: %v", patchName, appName, err)
	}

	for i, patchFunction := range patchFunctions {
		err := patchFunction(transaction)
		if err != nil {
			log.Fatal("Apps", "Cannot execute function %d to patch %s app %s: %v", i, patchName, appName, err)
		}
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		log.Fatal("Apps", "Cannot commit patch %s for app %s: %v", patchName, appName, err)
	}

	err = applyPatch(appName, patchName)
	if err != nil {
		log.Fatal("Apps", "Cannot register patch %s for app %s: %v", patchName, appName, err)
	}
}

// patchApplied checks if the patch is already applied.
func patchApplied(appName string, patchName string) bool {
	patch, _, err := client.NewClient().AppsApi.
		GetPatchByName(client.AuthenticationContext(), appName, patchName).
		Execute()
	if err != nil || !patch.Applied.IsSet() {
		return false
	}
	return *patch.Applied.Get()
}

// applyPatch marks that the patch is now applied.
func applyPatch(appName string, patchName string) error {
	_, err := client.NewClient().AppsApi.
		PatchPatchByName(client.AuthenticationContext(), appName, patchName).
		Apply(true).
		Execute()
	tools.LogError(err)
	return err
}
