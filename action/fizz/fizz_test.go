// This package extract and validate data
package fizz

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/manager"
	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"
	"github.com/arckadious/fizzbuzz/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFizz(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)
	cf, err := config.New("../../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)

	m := manager.New(cf, *response.New(200, "success", []response.ApiError{}, nil), validator.New(), &repository.Repository{})
	require.NotNil(m)
	mf := manager.NewFizz(m, repository.NewFizz(repository.New(&database.DB{})))
	require.NotNil(mf)

	//////////////////////
	// FizzAction.New() //
	//////////////////////

	fa := New(mf)
	require.NotNil(fa)

	///////////////////////////////////
	// FizzAction.HandleStatistics() //
	///////////////////////////////////

	// No tests needed, because it's just one call
	// to ac.mng.HandleStatistics() function, which is already tested

	/////////////////////////////
	// FizzAction.HandleFizz() //
	/////////////////////////////

	// body empty
	r := httptest.NewRequest("POST", "/v1/fizzbuzz", bytes.NewBufferString(""))
	w := httptest.NewRecorder()
	fa.HandleFizz(w, r)
	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"INVALID_DATA_ERR\",\"message\":\"Data JSON input bad format.\"}],\"data\":null}", w.Body.String())

	// data OK
	data := model.Input{Limit: 30, Multiples: []model.Multiple{
		{
			IntX: 5,
			StrX: "buzz",
		},
		{
			IntX: 3,
			StrX: "fizz",
		},
	}}
	buf, err := json.Marshal(data)
	require.NoError(err)
	r = httptest.NewRequest("POST", "/v1/fizzbuzz", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()
	fa.HandleFizz(w, r)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,buzzfizz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,buzzfizz\"}", w.Body.String())

	// data validatation error : Multiple array size unauthorized
	data.Multiples = append(data.Multiples, model.Multiple{IntX: 3, StrX: "fizz"})
	buf, err = json.Marshal(data)
	require.NoError(err)
	r = httptest.NewRequest("POST", "/v1/fizzbuzz", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()
	fa.HandleFizz(w, r)
	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"INVALID_FIELD_ERR\",\"message\":\"Key: 'Input.Multiples' Error:Field validation for 'Multiples' failed on the 'lte' tag\"}],\"data\":null}", w.Body.String())
}
