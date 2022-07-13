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
	"github.com/eliona-smart-building-assistant/go-eliona/api"
)

// UpsertHeap inserts or updates the given heap. If the heap with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertHeap(heap api.Heap) error {
	_, err := api.NewClient().HeapApi.PostHeap(context.Background()).Heap(heap).Execute()
	return err
}

// UpsertHeapIfAssetExists upsert the heap if the eliona id exists. Otherwise, the upsert is ignored
func UpsertHeapIfAssetExists[T any](heap api.Heap) error {
	exists, err := ExistAsset(heap.AssetId)
	if err != nil {
		return err
	}
	if exists {
		return UpsertHeap(heap)
	}
	return nil
}
