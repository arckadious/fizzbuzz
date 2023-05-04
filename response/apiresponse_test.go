// This package contains JSON response templates functions
package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiresponse(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)

	///////////////////////
	// ApiResponse.New() //
	///////////////////////

	// Message spécifié taille 1, data nil
	ar := New(http.StatusOK, "test", []ApiError{{Code: "test", Message: "test"}}, nil)
	require.NotNil(ar)
	w := httptest.NewRecorder()
	ar.WriteJSONResponse(w)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("{\"status\":\"test\",\"messages\":[{\"code\":\"test\",\"message\":\"test\"}],\"data\":null}", w.Body.String())

	//pas de message, data nil
	ar = New(http.StatusOK, "test", []ApiError{}, nil)
	assert.NotNil(ar)
	w = httptest.NewRecorder()
	ar.WriteJSONResponse(w)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("{\"status\":\"test\",\"messages\":[],\"data\":null}", w.Body.String())

	//message nil, data nil, setData = "test", setErrorResponse http.StatusBadRequest message 2
	ar = New(http.StatusOK, "test", []ApiError{}, nil)
	assert.NotNil(ar)
	w = httptest.NewRecorder()
	ar.SetData("test")
	ar.SetErrorResponse(http.StatusBadRequest, []ApiError{{Code: "test", Message: "test"}, {Code: "test2", Message: "test2"}})
	ar.WriteJSONResponse(w)
	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"test\",\"message\":\"test\"},{\"code\":\"test2\",\"message\":\"test2\"}],\"data\":\"test\"}", w.Body.String())

}
