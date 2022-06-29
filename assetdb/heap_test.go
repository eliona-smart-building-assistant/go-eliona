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

package assetdb

import (
	"context"
	"encoding/json"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"github.com/pashagolub/pgxmock"
	"testing"
	"time"
)

type Temperature struct {
	Value int
	Unit  string
}

func TestUpsertHeap(t *testing.T) {
	temperature := Temperature{35, "Celsius"}
	payload, _ := json.Marshal(temperature)
	mock := connectionMock()
	mock.ExpectExec("insert into heap").
		WithArgs(2, InfoSubtype, pgxmock.AnyArg(), payload).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	defer mock.Close(context.Background())
	_ = UpsertHeap(mock, Heap[Temperature]{2, InfoSubtype, time.Time{}, temperature})
}

func connectionMock() pgxmock.PgxConnIface {
	mock, err := pgxmock.NewConn()
	if err != nil {
		log.Fatal("database", "An error '%s' was not expected when opening a mocked database connection", err)
	}
	return mock
}
