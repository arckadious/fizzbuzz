// // This package creates and run a rest API server, using Gin framework
package server

// import (
// 	"database/sql"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/arckadious/fizzbuzz/config"
// 	cst "github.com/arckadious/fizzbuzz/constant"
// 	"github.com/arckadious/fizzbuzz/container"
// 	"github.com/arckadious/fizzbuzz/database"
// 	"github.com/arckadious/fizzbuzz/manager"
// 	"github.com/arckadious/fizzbuzz/model"
// 	"github.com/arckadious/fizzbuzz/response"
// 	"github.com/arckadious/fizzbuzz/server"
// 	"github.com/arckadious/fizzbuzz/validator"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// type RepositoryMock struct {
// 	mock.Mock
// }

// func (m *RepositoryMock) LogToDB(logType, msg, url, corID, checksum, status string) (err error) {
// 	args := m.Called()
// 	return args.Error(0)
// }

// func TestMiddleware(t *testing.T) {

// 	assert := assert.New(t)
// 	require := require.New(t)
// 	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
// 	require.NoError(err)

// 	repo := new(RepositoryMock)
// 	mng := manager.New(cf, *response.New(http.StatusOK, cst.StatusSuccess, make([]response.ApiError, 0), nil), validator.New(), repo)

// 	//////////////////
// 	// Server.New() //
// 	//////////////////
// 	// ctnr := container.Container{

// 	}
// 	mf := NewFizz(mng, repo)

// 	//////////////////////////////////////////
// 	// Manager Methods from Parent class OK //
// 	//////////////////////////////////////////

// 	mf.GetApiResponse()
// 	mf.GetValidator()

// 	///////////////////////
// 	// Fizz.HandleFizz() //
// 	///////////////////////

// 	// model.Input empty
// 	w := httptest.NewRecorder()
// 	mf.HandleFizz(w, model.Input{})
// 	require.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"\"}", w.Body.String())
// 	assert.Equal(http.StatusOK, w.Code)

// 	// limit 5 parameter
// 	w = httptest.NewRecorder()
// 	mf.HandleFizz(w, model.Input{Limit: 5})
// 	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,3,4,5\"}", w.Body.String())
// 	assert.Equal(http.StatusOK, w.Code)

// 	// Limit 30, one multiple 3
// 	w = httptest.NewRecorder()
// 	mf.HandleFizz(w, model.Input{Limit: 30, Multiples: []model.Multiple{
// 		{
// 			IntX: 3,
// 			StrX: "fizz",
// 		},
// 	}})
// 	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,5,fizz,7,8,fizz,10,11,fizz,13,14,fizz,16,17,fizz,19,20,fizz,22,23,fizz,25,26,fizz,28,29,fizz\"}", w.Body.String())
// 	assert.Equal(http.StatusOK, w.Code)

// 	// Limit 30, two multiple 3, and 5
// 	w = httptest.NewRecorder()
// 	mf.HandleFizz(w, model.Input{Limit: 30, Multiples: []model.Multiple{
// 		{
// 			IntX: 3,
// 			StrX: "fizz",
// 		},
// 		{
// 			IntX: 5,
// 			StrX: "buzz",
// 		},
// 	}})

// 	assert.Equal(http.StatusOK, w.Code)
// 	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,fizzbuzz\"}", w.Body.String())

// 	// Limit 30, two multiple 5, and 3
// 	w = httptest.NewRecorder()
// 	mf.HandleFizz(w, model.Input{Limit: 30, Multiples: []model.Multiple{
// 		{
// 			IntX: 5,
// 			StrX: "buzz",
// 		},
// 		{
// 			IntX: 3,
// 			StrX: "fizz",
// 		},
// 	}})
// 	assert.Equal(http.StatusOK, w.Code)
// 	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,buzzfizz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,buzzfizz\"}", w.Body.String())

// 	/////////////////////////////
// 	// Fizz.HandleStatistics() //
// 	/////////////////////////////

// 	// No rows in database (mock DB)
// 	repo.On("GetMostRequestUsed").Return("", 0, true, sql.ErrNoRows)
// 	w = httptest.NewRecorder()
// 	mf.HandleStatistics(w)
// 	repo.AssertExpectations(t)
// 	assert.Equal(http.StatusPartialContent, w.Code)
// 	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":null}", w.Body.String())

// 	// Wrong message returned by function GetMostRequestUsed.
// 	repo.ExpectedCalls = []*mock.Call{}
// 	repo.On("GetMostRequestUsed").Return("toto", 0, false, nil)
// 	w = httptest.NewRecorder()
// 	mf.HandleStatistics(w)
// 	repo.AssertExpectations(t)
// 	assert.Equal(http.StatusInternalServerError, w.Code)
// 	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"invalid character 'o' in literal true (expecting 'r')\"}],\"data\":null}", w.Body.String())

// 	// Wrong JSON message returned by function GetMostRequestUsed.
// 	repo.ExpectedCalls = []*mock.Call{}
// 	repo.On("GetMostRequestUsed").Return("{\"toto\":\"test\"}", 0, false, nil)
// 	w = httptest.NewRecorder()
// 	mf.HandleStatistics(w)
// 	repo.AssertExpectations(t)
// 	assert.Equal(http.StatusInternalServerError, w.Code)
// 	assert.Contains(w.Body.String(), "{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"Key: ")

// 	// Process OK
// 	repo.ExpectedCalls = []*mock.Call{}
// 	repo.On("GetMostRequestUsed").Return(`{
// 			"multiples": [
// 			  {
// 				"intX": 3,
// 				"strX": "fizz"
// 			  },
// 			  {
// 				"intX": 5,
// 				"strX": "buzz"
// 			  }
// 			],
// 			"limit": 100
// 		  }`, 56, false, nil)
// 	w = httptest.NewRecorder()
// 	mf.HandleStatistics(w)
// 	repo.AssertExpectations(t)
// 	assert.Equal(http.StatusOK, w.Code)
// 	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":{\"request\":{\"Limit\":100,\"Multiples\":[{\"IntX\":3,\"StrX\":\"fizz\"},{\"IntX\":5,\"StrX\":\"buzz\"}]},\"hits\":56}}", w.Body.String())

// }
