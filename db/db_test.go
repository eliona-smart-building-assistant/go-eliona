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

package db

import (
	"fmt"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Temperature struct {
	Value int
	Unit  string
}

type Weather struct {
	Temperature *Temperature
	Remark      *string
	DayOfWeek   *int
}

func TestPointer(t *testing.T) {
	mock := connectionMock()
	rows := mock.NewRows([]string{"temperature", "remark", "day_of_week"}).
		AddRow(nil, "Cloudy", 4)
	mock.ExpectQuery("select (.+) from weather").WillReturnRows(rows)
	weather, err := QuerySingleRow[Weather](mock, "select value, unit from weather")
	assert.Nil(t, err)
	assert.Equal(t, "", weather.Temperature.Unit)
	assert.Equal(t, "Cloudy", *weather.Remark)
	assert.Equal(t, 4, *weather.DayOfWeek)
}

func TestQuery(t *testing.T) {
	mock := connectionMock()
	rows := mock.NewRows([]string{"value", "unit"}).
		AddRow(25, "Celsius").
		AddRow(67, "Fah").
		AddRow(40, "Celsius")
	mock.ExpectQuery("select (.+) from temperatures").WillReturnRows(rows)

	resultsChan := make(chan Temperature)
	go func() {
		err := Query(mock, "select value, unit from temperatures", resultsChan)
		if err != nil {
			assert.Nil(t, err)
		}
	}()

	var results []Temperature
	for result := range resultsChan {
		results = append(results, result)
	}

	assert.Equal(t, results[0], Temperature{25, "Celsius"})
	assert.Equal(t, results[1], Temperature{67, "Fah"})
	assert.Equal(t, results[2], Temperature{40, "Celsius"})
}

func TestQuerySingleRowError(t *testing.T) {
	mock := connectionMock()
	mock.ExpectQuery("select.+from temperatures").WillReturnError(fmt.Errorf("error"))
	result, err := QuerySingleRow[Temperature](mock, "select value, unit from temperatures")
	assert.Equal(t, fmt.Errorf("error"), err)
	assert.Equal(t, result.Value, 0)
}

func TestQuerySingleRow(t *testing.T) {
	mock := connectionMock()
	rows := mock.NewRows([]string{"value", "unit"}).
		AddRow(25, "Celsius")
	mock.ExpectQuery("select (.+) from temperatures").WillReturnRows(rows)
	result, err := QuerySingleRow[Temperature](mock, "select value, unit from temperatures")
	assert.Nil(t, err)
	assert.Equal(t, result, Temperature{25, "Celsius"})
}

func TestQueryError(t *testing.T) {
	mock := connectionMock()
	mock.ExpectQuery("select.+from temperatures").WillReturnError(fmt.Errorf("error"))
	resultsChan := make(chan Temperature)
	go func() {
		err := Query(mock, "select value, unit from temperatures", resultsChan)
		if err != nil {
			assert.Equal(t, fmt.Errorf("error"), err)
		}
	}()
	assert.Equal(t, (<-resultsChan).Value, 0)
}

func TestExecuteFile(t *testing.T) {
	mock := connectionMock()
	mock.ExpectExec("create table test").WillReturnResult(pgxmock.NewResult("CREATE", 0))
	mock.ExpectExec("insert into test").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	err := ExecFile(mock, "test.sql")
	assert.Nil(t, err)
}

func TestInsert(t *testing.T) {
	mock := connectionMock()
	mock.ExpectExec("insert into temperatures").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := Exec(mock, "insert into temperatures (value, unit) values (?, ?)", 25, "Celsius")
	assert.Nil(t, err)
}

func connectionMock() pgxmock.PgxConnIface {
	mock, err := pgxmock.NewConn()
	if err != nil {
		log.Fatal("database", "An error '%s' was not expected when opening a mocked database connection", err)
	}
	return mock
}
