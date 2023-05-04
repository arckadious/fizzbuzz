// This package process data from action class handlers
package manager

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	cst "github.com/arckadious/fizzbuzz/constant"
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
	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)
	mng := New(cf, *response.New(http.StatusOK, cst.StatusSuccess, make([]response.ApiError, 0), nil), validator.New(), &repository.Repository{})
	// repoFizz := new(FizzRepositoryMock)

	////////////////
	// Fizz.New() //
	////////////////

	mf := NewFizz(mng, &repository.Fizz{}) //repoFizz)

	//////////////////////////////////////////
	// Manager Methods from Parent class OK //
	//////////////////////////////////////////

	mf.GetApiResponse()
	mf.GetValidator()

	///////////////////////
	// Fizz.HandleFizz() //
	///////////////////////

	// model.Input empty
	w := httptest.NewRecorder()
	mf.HandleFizz(w, model.Input{})
	require.Equal(w.Body.String(), "{\"status\":\"success\",\"messages\":[],\"data\":\"\"}")
	assert.Equal(200, w.Code)

	// limit 5 parameter
	w = httptest.NewRecorder()
	mf.HandleFizz(w, model.Input{Limit: 5})
	assert.Equal(w.Body.String(), "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,3,4,5\"}")
	assert.Equal(200, w.Code)

	// Limit 30, one multiple 3
	w = httptest.NewRecorder()
	mf.HandleFizz(w, model.Input{Limit: 30, Multiples: []model.Multiple{
		{
			IntX: 3,
			StrX: "fizz",
		},
	}})
	assert.Equal(w.Body.String(), "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,5,fizz,7,8,fizz,10,11,fizz,13,14,fizz,16,17,fizz,19,20,fizz,22,23,fizz,25,26,fizz,28,29,fizz\"}")
	assert.Equal(200, w.Code)

	// Limit 30, two multiple 3, and 5
	w = httptest.NewRecorder()
	mf.HandleFizz(w, model.Input{Limit: 30, Multiples: []model.Multiple{
		{
			IntX: 3,
			StrX: "fizz",
		},
		{
			IntX: 5,
			StrX: "buzz",
		},
	}})

	assert.Equal(200, w.Code)
	assert.Equal(w.Body.String(), "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,fizzbuzz\"}")

	// Limit 30, two multiple 5, and 3
	w = httptest.NewRecorder()
	mf.HandleFizz(w, model.Input{Limit: 30, Multiples: []model.Multiple{
		{
			IntX: 5,
			StrX: "buzz",
		},
		{
			IntX: 3,
			StrX: "fizz",
		},
	}})
	assert.Equal(200, w.Code)
	assert.Equal(w.Body.String(), "{\"status\":\"success\",\"messages\":[],\"data\":\"1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,buzzfizz,16,17,fizz,19,buzz,fizz,22,23,fizz,buzz,26,fizz,28,29,buzzfizz\"}")

	/////////////////////////////
	// Fizz.HandleStatistics() //
	/////////////////////////////

	// // No rows in database (mocked)
	// repoFizz.On("GetMostRequestUsed").Return("", 0, true, nil)

	// // call the code we are testing
	// mf.HandleStatistics(w)

	// // assert that the expectations were met
	// repoFizz.AssertExpectations(t)
}
