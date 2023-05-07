// This package contains various useful tools like MD5 Hash and http body extract
package util

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {

	assert := assert.New(t)

	///////////////////
	// GenerateUID() //
	///////////////////

	uuid, err := GenerateUID()
	assert.NoError(err)
	assert.NotEqual("", uuid)

	///////////////////
	// ExtractBody() //
	///////////////////

	body, err := ExtractBody(nil) //request nil
	assert.Error(err)
	assert.Equal("", string(body))

	r, _ := http.NewRequest("GET", "http://google.fr", bytes.NewBufferString("test")) //body already read
	_, err = io.ReadAll(r.Body)
	assert.NoError(err)
	body, err = ExtractBody(r)
	assert.NoError(err)
	assert.Equal("", string(body))

	r, _ = http.NewRequest("GET", "http://google.fr", nil) //body nil
	body, err = ExtractBody(r)
	assert.NoError(err)
	assert.Equal("", string(body))

	r, _ = http.NewRequest("GET", "http://google.fr", iotest.ErrReader(errors.New("test"))) //body error
	body, err = ExtractBody(r)
	assert.Error(err)
	assert.Equal("", string(body))

	r, _ = http.NewRequest("GET", "http://google.fr", bytes.NewBufferString("test")) //OK : body extracted and still present on request
	body, err = ExtractBody(r)
	assert.NoError(err)
	assert.Equal("test", string(body))
	rBody, err := io.ReadAll(r.Body)
	assert.NoError(err)
	assert.Equal("test", string(rBody))

	//////////////////
	// GetMD5Hash() //
	//////////////////

	assert.Equal("d41d8cd98f00b204e9800998ecf8427e", GetMD5Hash(""))
	assert.Equal("098f6bcd4621d373cade4e832627b4f6", GetMD5Hash("test"))

}
