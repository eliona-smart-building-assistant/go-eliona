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

package eliona

import (
	"fmt"
	"github.com/eliona-smart-building-assistant/go-eliona/api"
	"os"
	"testing"
	"time"
)

type Temperature struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

func TestUpsertHeap(t *testing.T) {
	os.Setenv("API_ENDPOINT", "http://localhost:8888/apps/v2")
	temperature := Temperature{35, "Celsius"}
	err := UpsertHeap(api.Heap{AssetId: 2, Subtype: api.INPUT, Timestamp: &time.Time{}, Data: api.StructToMap(temperature)})
	if err != nil {
		fmt.Println(err)
	}
}
