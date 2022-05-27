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

package assets

import (
	"encoding/json"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"github.com/jackc/pgx/v4"
	"time"
)

// Heap defines the heap data which specific payload data of type T
type Heap[T any] struct {
	AssetId   int       `json:"asset_id"`
	Subtype   Subtype   `json:"subtype"`
	TimeStamp time.Time `json:"ts"`
	Data      T         `json:"data"`
}

func (heap Heap[T]) GetData() T {
	return heap.Data
}

// UpsertHeap inserts or updates the given heap. If the heap with the specified subtype does not exists, it will be created.
// Otherwise, the timestamp and the data are updated.
func UpsertHeap[T any](connection db.Connection, heap Heap[T]) error {
	payload, err := json.Marshal(heap.Data)
	if err != nil {
		log.Error("Assets", "Failed to marshal value: %s", err.Error())
		return err
	}
	return db.Exec(connection,
		"insert into heap (asset_id, subtype, ts, data) "+
			"values ($1, $2, $3, $4) "+
			"on conflict (asset_id, subtype) "+
			"do update set ts = excluded.ts, data = excluded.data",
		heap.AssetId, string(heap.Subtype), heap.TimeStamp, payload)
}

// ListenHeap listens on the channel heap and collect updated or inserted heaps
func ListenHeap[T any](connection *pgx.Conn, outputs chan Heap[T]) {
	db.Listen(connection, "heap", outputs)
}
