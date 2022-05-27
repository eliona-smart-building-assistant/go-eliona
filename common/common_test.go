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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetenv(t *testing.T) {
	assert.Equal(t, "default", Getenv("FOO", "default"))
	t.Setenv("FOO", "bar")
	assert.Equal(t, "bar", Getenv("FOO", "default"))
	t.Setenv("FOO", "")
	assert.Equal(t, "", Getenv("FOO", "default"))
}

func TestAppName(t *testing.T) {
	assert.Equal(t, "", AppName())
	t.Setenv("APPNAME", "foobar")
	assert.Equal(t, "foobar", AppName())
}
