// This package process data from action class handlers
package manager

import (
	"net/http"
	"testing"

	"github.com/arckadious/fizzbuzz/config"
	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)
	cf, err := config.New("../tests/mock/parametersOK.json", *validator.New())
	require.NoError(err)

	///////////////////
	// Manager.New() //
	///////////////////

	vexpected := validator.New()
	mng := New(cf, *response.New(http.StatusOK, cst.StatusSuccess, make([]response.ApiError, 0), nil), vexpected, &repository.Repository{})

	//////////////////////////////
	// Manager.GetApiResponse() //
	//////////////////////////////

	apiResponse := mng.GetApiResponse()
	assert.Equal(response.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     cst.StatusSuccess,
		Messages:   make([]response.ApiError, 0),
	}, apiResponse)

	//////////////////////////////
	// Manager.GetValidator() //
	//////////////////////////////
	assert.Equal(vexpected, mng.GetValidator())
}
