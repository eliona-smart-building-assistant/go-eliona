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

package asset

import (
	"context"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
)

// UpsertData inserts or updates the given heap. If the heap with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertData(heap api.Data) error {
	_, err := client.NewClient().DataApi.PutData(context.Background()).Data(heap).Execute()
	return err
}

// UpsertDataIfAssetExists upsert the heap if the eliona id exists. Otherwise, the upsert is ignored
func UpsertDataIfAssetExists[T any](heap api.Data) error {
	exists, err := ExistAsset(heap.AssetId)
	if err != nil {
		return err
	}
	if exists {
		return UpsertData(heap)
	}
	return nil
}
