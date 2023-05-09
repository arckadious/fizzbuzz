// This package creates and run a rest API server, using Gin framework
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

// Init server for all package tests just once.
func InitServer(t *testing.T) (*assert.Assertions, *require.Assertions, *test.Hook) {
	assert := assert.New(t)
	require := require.New(t)
	hook := new(test.Hook)
	logrus.AddHook(hook)

	if sG != nil {
		return assert, require, hook
	}

	logrus.SetLevel(logrus.ErrorLevel)
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

// Test boot server error
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
		var fatalEntry logrus.Entry
		isFatal := false
		for _, entry := range hook.Entries {
			if entry.Level == logrus.FatalLevel {
				fatalEntry = entry
				isFatal = true
			}
		}
		if assert.True(isFatal) {
			assert.Contains(fatalEntry.Message, "address already in use")
		}
	}
	logrus.StandardLogger().ExitFunc = logrus.New().ExitFunc

}

// Test unexpected cases handlers
func TestServerHandlers(t *testing.T) {
	assert, _, _ := InitServer(t)

	router := sG.handler()

	// recoveryHandler
	w := httptest.NewRecorder()
	ctx := gin.CreateTestContextOnly(w, router)
	sG.recoveryHandler(ctx, errors.New("test"))
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

// Test endpoints and Logger middleware if database is closed (database unavailable approach, we only want to see errors)
func TestServerDBUnavailable(t *testing.T) {

	assert, require, hook := InitServer(t)

	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	cf.Port = 8008
	require.NoError(err)

	validator := validator.New()
	s := New( // Use a new server for database down cases
		container.New(
			cf,
			validator,
			database.New(cf),
		),
	)

	router := s.handler()

	// Endpoint /v1/statistics database closed
	s.container.Db.Shutdown()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/statistics", nil)
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"INTERNAL_SERVER_ERR\",\"message\":\"sql: database is closed\"}],\"data\":null}", w.Body.String())

	//Endpoint /v1/fizzbuzz database closed
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
		"limit": 30
	  }`))
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,fizzbuzz\"}", w.Body.String())

	//check last error log (logger works under go routines)
	if assert.Eventually(func() bool { return hook.LastEntry() != nil }, 5*time.Second, 10*time.Millisecond) {
		assert.Equal("Logger coudn't send response data: sql: database is closed", hook.LastEntry().Message)
	}
}

// Simple call to endpoints
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
		"limit": 30
	  }`))
	router.ServeHTTP(w, req)
	assert.Equal(200, w.Code)
	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,fizzbuzz\"}", w.Body.String())
}
