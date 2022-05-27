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

package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

const test = "test"

// go test -race dev.azure.com/itec-ag/eliona/edge-node/services/ttn-client/log
func TestLogger_SetLevelRace(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.SetLevel(DebugLevel)
	for i := 0; i < 100; i++ {
		i := i
		go func() {
			// changing levels in multiple goroutines
			l.SetLevel(Level(i % 4))
			l.Println(Level(i%4), "TEST", "testmessage", test)
		}()
		l.Println(Level(i%4), "TEST", "testmessage")
	}
}

func TestSetLevel(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.SetLevel(DebugLevel)

	levels := []Level{ErrorLevel, WarnLevel, InfoLevel, DebugLevel}
	for i, k := range levels {
		l.SetLevel(k)
		for _, kk := range levels {
			l.Println(kk, "TST", "test")
		}
		if n := strings.Count(b.String(), "\n"); n != i+1 {
			t.Errorf("expected %d lines, got %d", i+1, n)
		}
		b.Reset()
	}
}

func TestHandyMethods(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.SetLevel(DebugLevel)
	l.Info("foobar", "foo %s", "bar")
	assert.Contains(t, b.String(), "INFO")
	assert.Contains(t, b.String(), "\tfoobar\tfoo bar\n")
}

func TestLevelSettingViaEnvironment(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	assert.Equal(t, InfoLevel, l.level())

	t.Setenv("LOG_LEVEL", "error")
	l = New(&b)
	assert.Equal(t, ErrorLevel, l.level())
	l.Info("foobar", "foo %s", "bar")
	assert.NotContains(t, b.String(), "INFO")

	t.Setenv("LOG_LEVEL", "debug")
	l = New(&b)
	assert.Equal(t, DebugLevel, l.level())
	l.Info("foobar", "foo %s", "bar")
	assert.Contains(t, b.String(), "INFO")
}

func TestSetBufLimit(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.SetLevel(DebugLevel)
	bigString := strings.Repeat("testing length", 100)
	// by default buffer is not limited
	l.Println(DebugLevel, "TST", bigString)
	if !strings.HasSuffix(b.String(), bigString+"\n") {
		t.Errorf("log output should match\n%q\nis\n%q", b.String(), bigString)
	}
	b.Reset()

	buflimit := int64(1024)
	l.SetBufLimit(buflimit)
	l.Println(DebugLevel, "TST", bigString)
	if int64(b.Len()) != buflimit {
		t.Errorf("log length should be equal to buf limit %d is %d", b.Len(), buflimit)
	}
	if b.Bytes()[b.Len()-1] != '\n' {
		t.Errorf("last byte should be equal \\n")
	}
}

func TestOutput(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.SetLevel(DebugLevel)
	l.Println(DebugLevel, "TST", "test")
	re := regexp.MustCompile("DEBUG\\t\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}.\\d{6}\\tTST\\ttest\\n")
	if !re.Match(b.Bytes()) {
		t.Errorf("log output should match pattern %q is %q", re.String(), b.String())
	}
}

func TestEmptyPrintCreatesLine(t *testing.T) {
	var b bytes.Buffer
	l := New(&b)
	l.SetLevel(DebugLevel)
	l.Print(DebugLevel, "Header", "")
	l.Println(DebugLevel, "Header", "non-empty")
	output := b.String()
	if n := strings.Count(output, "Header"); n != 2 {
		t.Errorf("expected 2 headers, got %d", n)
	}
	if n := strings.Count(output, "\n"); n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func BenchmarkPrintln(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(&buf)
	l.SetLevel(DebugLevel)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(DebugLevel, "TST", testString)
	}
}
