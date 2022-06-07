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

package apps

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunOnce(t *testing.T) {
	var counter int64 = 1
	runsOnlyOnceAtSameTime(&counter)
	assert.Equal(t, int64(2), counter)
	RunOnce(func() { runsOnlyOnceAtSameTime(&counter) }, 4711)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, int64(3), counter)
	RunOnce(func() { runsOnlyOnceAtSameTime(&counter) }, 4711)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, int64(3), counter)
	RunOnce(func() { runsOnlyOnceAtSameTime(&counter) }, 4712)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, int64(4), counter)
	RunOnce(func() { runsOnlyOnceAtSameTime(&counter) }, 4711)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, int64(4), counter)
	time.Sleep(time.Millisecond * 100)
	RunOnce(func() { runsOnlyOnceAtSameTime(&counter) }, 4711)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, int64(5), counter)
}

func runsOnlyOnceAtSameTime(counter *int64) {
	atomic.AddInt64(counter, 1)
	time.Sleep(time.Millisecond * 100)
}
