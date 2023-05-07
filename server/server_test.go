// This package creates and run a rest API server, using Gin framework
package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/container"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/response"
	"github.com/arckadious/fizzbuzz/validator"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sG *Server

func InitServer(t *testing.T) (*assert.Assertions, *require.Assertions, *test.Hook) {
	assert := assert.New(t)
	require := require.New(t)
	hook := new(test.Hook)
	logrus.AddHook(hook)
	if sG != nil {
		return assert, require, hook
	}

	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)

	validator := validator.New()
	sG = New(
		container.New(
			cf,
			validator,
			database.New(cf),
		),
	)

	return assert, require, hook
}

func TestServerErrorBindPortAlreadyUsed(t *testing.T) {

	assert, _, hook := InitServer(t)
	var fatal bool
	go sG.Run()

	ch := make(chan struct{}, 1)

	logrus.StandardLogger().ExitFunc = func(int) {
		fatal = true
		ch <- struct{}{}
	}

	go sG.Run()
	<-ch

	assert.True(fatal)
	if assert.NotNil(hook.LastEntry()) {
		assert.Contains(hook.LastEntry().Message, "bind: address already in use")
	}
	logrus.StandardLogger().ExitFunc = logrus.New().ExitFunc

}

func TestServerHandlers(t *testing.T) {
	assert, _, _ := InitServer(t)

	router := sG.handler()

	// recoveryHandler
	w := httptest.NewRecorder()
	ctx := gin.CreateTestContextOnly(w, router)
	sG.recoveryHandler(ctx, nil)
	assert.Equal(500, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"Internal Server Error, oups !\"}],\"data\":null}", w.Body.String())

	// notFoundHandler
	w = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/unkwown", nil)
	router.ServeHTTP(w, req)
	assert.Equal(404, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"ROUTE_NOT_FOUND_ERR\",\"message\":\"This is not what you are looking for.\"}],\"data\":null}", w.Body.String())

	// methodNotAllowedHandler
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/ping", nil)
	router.ServeHTTP(w, req)
	assert.Equal(405, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"METHOD_NOT_ALLOWED_ERR\",\"message\":\"Method is not allowed, boy.\"}],\"data\":null}", w.Body.String())

}

func TestServerEndpoints(t *testing.T) {
	assert, _, _ := InitServer(t)

	router := sG.handler()

	// Endpoint /ping
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.Equal("Ping OK !", w.Body.String())

	// Endpoint /swagger
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/swagger", nil)
	router.ServeHTTP(w, req)
	assert.Equal(301, w.Code)
	assert.Equal("<a href=\"/swagger/\">Moved Permanently</a>.\n\n", w.Body.String())

	// Endpoint /v1/statistics
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/statistics", nil)
	router.ServeHTTP(w, req)
	assert.True(w.Code == 206 || w.Code == 200)
	assert.NoError(json.Unmarshal(w.Body.Bytes(), &response.ApiResponse{}))

	//Endpoint /v1/fizzbuzz
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/v1/fizzbuzz", bytes.NewBufferString(`{
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
	  }`))
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.NoError(json.Unmarshal(w.Body.Bytes(), &response.ApiResponse{}))
}

// func TestServer(t *testing.T) {

// 	assert, require, hook := InitServer(t)
// 	sG.Run()

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
