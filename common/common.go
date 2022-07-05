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

package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// Map to check if a function started with RunOnce is currently running.
var runOnceIds sync.Map

// RunOnce starts a function if this function not currently running. RunOnce knows, with function is currently
// running (identified by id) and skips starting the function again.
func RunOnce(function func(), id any) {
	go func() {
		_, alreadyRuns := runOnceIds.Load(id)
		if !alreadyRuns {
			runOnceIds.Store(id, nil)
			function()
			runOnceIds.Delete(id)
		}
	}()
}

// WaitFor helps to start multiple functions in parallel and waits until all functions are completed. Normally one app has
// only one main functions that runs in an infinite loop, except the app is stopped externally, e.g. during a shut-down
// of the eliona environment.
func WaitFor(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))
	for _, function := range functions {
		function := function
		go func() {
			function()
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}

// Loop wraps a function in an endless loop and calls the function in the defined interval.
func Loop(function func(), interval time.Duration) func() {
	return func() {
		osSignals := make(chan os.Signal, 1)
		defer close(osSignals)
		signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		for {
			function()
			select {
			case <-time.After(interval):
			case <-osSignals:
				return
			}
		}
	}
}

// Getenv reads the value from environment variable named by key.
// If the key is not defined as environment variable the default string is returned.
func Getenv(key, fallback string) string {
	value, present := os.LookupEnv(key)
	if present {
		return value
	}
	return fallback
}

// UnmarshalFile returns the content of the file as object of type T
func UnmarshalFile[T any](path string) (T, error) {
	var object T
	data, err := ioutil.ReadFile(filepath.Join(path))
	if err != nil {
		return object, err
	}
	err = json.Unmarshal(data, &object)
	if err != nil {
		return object, err
	}
	return object, nil
}

// AppName returns the name of the app uses the library. The app name is defined in the
// environment variable APPNAME. If not defined, the AppName returns nil.
func AppName() string {
	return Getenv("APPNAME", "")
}

// Ptr delivers the pointer of any constant value like Ptr("foo")
func Ptr[T any](v T) *T {
	return &v
}

// StructToMap converts a struct to map of struct properties
func StructToMap(data any) map[string]interface{} {
	d, _ := json.Marshal(&data)
	var m map[string]interface{}
	_ = json.Unmarshal(d, &m)
	return m
}
