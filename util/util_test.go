package util

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestUtil(t *testing.T) {

	//GenerateUUID
	uuid, err := GenerateUUID()
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", uuid)

	//ExtractBody
	body, err := ExtractBody(nil) //request nil
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", string(body))

	r, _ := http.NewRequest("GET", "http://google.fr", bytes.NewBufferString("test")) //body already read
	io.ReadAll(r.Body)
	body, err = ExtractBody(r)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", string(body))

	r, _ = http.NewRequest("GET", "http://google.fr", errReader(0)) //body error
	body, err = ExtractBody(r)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", string(body))

	r, _ = http.NewRequest("GET", "http://google.fr", bytes.NewBufferString("test")) //OK : body extracted and still present on request
	body, err = ExtractBody(r)
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", string(body))
	rBody, err := io.ReadAll(r.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", string(rBody))

	//GetMD5Hash
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", GetMD5Hash(""))
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", GetMD5Hash("test"))

}
