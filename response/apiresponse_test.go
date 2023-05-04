// This package contains JSON response templates functions
package response

import (
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
	ar := New(200, "test", []ApiError{{Code: "test", Message: "test"}}, nil)
	require.NotNil(ar)
	w := httptest.NewRecorder()
	ar.WriteJSONResponse(w)
	assert.Equal(200, w.Code)
	assert.Equal("{\"status\":\"test\",\"messages\":[{\"code\":\"test\",\"message\":\"test\"}],\"data\":null}", w.Body.String())

	//pas de message, data nil
	ar = New(200, "test", []ApiError{}, nil)
	assert.NotNil(ar)
	w = httptest.NewRecorder()
	ar.WriteJSONResponse(w)
	assert.Equal(200, w.Code)
	assert.Equal("{\"status\":\"test\",\"messages\":[],\"data\":null}", w.Body.String())

	//message nil, data nil, setData = "test", setErrorResponse 400 message 2
	ar = New(200, "test", []ApiError{}, nil)
	assert.NotNil(ar)
	w = httptest.NewRecorder()
	ar.SetData("test")
	ar.SetErrorResponse(400, []ApiError{{Code: "test", Message: "test"}, {Code: "test2", Message: "test2"}})
	ar.WriteJSONResponse(w)
	assert.Equal(400, w.Code)
	assert.Equal("{\"status\":\"error\",\"messages\":[{\"code\":\"test\",\"message\":\"test\"},{\"code\":\"test2\",\"message\":\"test2\"}],\"data\":\"test\"}", w.Body.String())

}
