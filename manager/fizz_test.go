// This package process data from action class handlers
package manager

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"
	tests "github.com/arckadious/fizzbuzz/tests/mock"
	"github.com/arckadious/fizzbuzz/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var repoFizz *tests.FizzRepositoryMock
var mf *Fizz

func InitRepoAndManager(t *testing.T) (*assert.Assertions, *require.Assertions) {
	t.Helper()
	assert := assert.New(t)
	require := require.New(t)

	if repoFizz != nil && mf != nil {
		return assert, require
	}

	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)
	mng := New(cf, *response.New(http.StatusOK, cst.StatusSuccess, make([]response.ApiError, 0), nil), validator.New(), &repository.Repository{})
	repoFizz = new(tests.FizzRepositoryMock) //repository Fizz mock

	////////////////
	// Fizz.New() //
	////////////////

	mf = NewFizz(mng, repoFizz)

	return assert, require
}

// ///////////////////////////////////////
// Manager Methods from Parent class OK //
// ///////////////////////////////////////
func TestFizzParentMethods(t *testing.T) {

	InitRepoAndManager(t)

	mf.GetApiResponse()
	mf.GetValidator()
}

// ////////////////////
// Fizz.HandleFizz() //
// ////////////////////
func TestHandleFizz(t *testing.T) {

	assert, _ := InitRepoAndManager(t)

	var tests = []struct {
		name     string
		i        model.Input
		wantCode int
		wantBody string
	}{
		{"model.Input empty", model.Input{}, http.StatusOK, "{\"status\":\"success\",\"messages\":[],\"data\":\"\"}"},
		{"limit 5 parameter", model.Input{Limit: 5}, http.StatusOK, "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,3,4,5\"}"},

		{"Limit 30, one multiple 3", model.Input{Limit: 30, Multiples: []model.Multiple{
			{
				IntX: 3,
				StrX: "fizz",
			},
		}}, http.StatusOK, "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,5,fizz,7,8,fizz,10,11,fizz,13,14,fizz,16,17,fizz,19,20,fizz,22,23,fizz,25,26,fizz,28,29,fizz\"}"},

		{"Limit 30, two multiple 3, and 5", model.Input{Limit: 30, Multiples: []model.Multiple{
			{
				IntX: 3,
				StrX: "fizz",
			},
			{
				IntX: 5,
				StrX: "buzz",
			},
		}}, http.StatusOK, "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,fizzbuzz\"}"},

		{"Limit 30, two multiple 5, and 3", model.Input{Limit: 30, Multiples: []model.Multiple{
			{
				IntX: 5,
				StrX: "buzz",
			},
			{
				IntX: 3,
				StrX: "fizz",
			},
		}}, http.StatusOK, "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,buzzfizz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,buzzfizz\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			mf.HandleFizz(w, tt.i)
			assert.Equal(tt.wantBody, w.Body.String(), tt.name)
			assert.Equal(tt.wantCode, w.Code, tt.name)
		})
	}
}

// //////////////////////////
// Fizz.HandleStatistics() //
// //////////////////////////
func TestHandleStatistics(t *testing.T) {

	// Mock repository is used to test HandleStatistics function.
	assert, _ := InitRepoAndManager(t)

	var tests = []struct {
		name     string
		contains bool
		i        []interface{}
		wantCode int
		wantBody string
	}{
		{"No rows in database (mock DB)", false, []interface{}{"", 0, true, sql.ErrNoRows}, http.StatusPartialContent, "{\"status\":\"success\",\"messages\":[],\"data\":null}"},
		{"Error from DB", false, []interface{}{"", 0, false, errors.New("test")}, http.StatusInternalServerError, "{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"test\"}],\"data\":null}"},
		{"Wrong message returned by function GetMostRequestUsed", false, []interface{}{"toto", 0, false, nil}, http.StatusInternalServerError, "{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"invalid character 'o' in literal true (expecting 'r')\"}],\"data\":null}"},
		{"Wrong JSON message returned by function GetMostRequestUsed", true, []interface{}{"{\"toto\":\"test\"}", 0, false, nil}, http.StatusInternalServerError, "{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"Key: "},
		{"Process OK", false, []interface{}{`{
			"multiples": [
			  {
				"intX": 3,
				"strX": "fizz"
			  },
			  {
				"intX": 5,
				"strX": "buzz"
			  }
			],
			"limit": 100
		  }`, 56, false, nil}, http.StatusOK, "{\"status\":\"success\",\"messages\":[],\"data\":{\"request\":{\"Limit\":100,\"Multiples\":[{\"IntX\":3,\"StrX\":\"fizz\"},{\"IntX\":5,\"StrX\":\"buzz\"}]},\"hits\":56}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoFizz.ExpectedCalls = []*mock.Call{}
			repoFizz.On("GetMostRequestUsed").Return(tt.i...)
			w := httptest.NewRecorder()
			mf.HandleStatistics(w)
			repoFizz.AssertExpectations(t)
			assert.Equal(tt.wantCode, w.Code, tt.name)
			if tt.contains {
				assert.Contains(w.Body.String(), tt.wantBody, tt.name)
			} else {
				assert.Equal(tt.wantBody, w.Body.String(), tt.name)
			}
		})
	}
}
