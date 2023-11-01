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

package dashboard

import (
	"fmt"
	"path/filepath"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona-api-client/v2/tools"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

// UpsertWidgetType insert or updates an asset and returns the id
func UpsertWidgetType(widgetType api.WidgetType) error {
	_, _, err := client.NewClient().WidgetsTypesAPI.
		PutWidgetType(client.AuthenticationContext()).
		Expansions([]string{"WidgetType.elements"}).
		WidgetType(widgetType).
		Execute()
	tools.LogError(err)
	return err
}

func initWidgetTypeFile(path string) error {
	widgetType, err := common.UnmarshalFile[api.WidgetType](path)
	if err != nil {
		return fmt.Errorf("unmarshalling file %s: %v", path, err)
	}
	return UpsertWidgetType(widgetType)
}

// InitWidgetTypeFile inserts or updates the type build from the content of the given file.
func InitWidgetTypeFile(path string) func(db.Connection) error {
	return func(db.Connection) error {
		return initWidgetTypeFile(path)
	}
}

func InitWidgetTypeFiles(pattern string) func(db.Connection) error {
	return func(db.Connection) error {
		paths, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("glob file pattern %s: %v", pattern, err)
		}
		for _, path := range paths {
			err := initWidgetTypeFile(path)
			if err != nil {
				return fmt.Errorf("initializing widget type %s: %v", path, err)
			}
		}
		return nil
	}
}
