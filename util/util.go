// This package contains various useful tools like MD5 Hash and http body extract
package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GenerateUID creates a unique ID (used to make a corelation between http requests and responses in database)
func GenerateUID() (uuid string, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return
	} else {
		uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	}
	return
}

// ExtractBody extracts a copy of http request body
func ExtractBody(r *http.Request) (body []byte, err error) {
	if r == nil {
		err = errors.New("Extract body : request nil")
		return
	}

	if r.Body == nil {
		return
	}

	body, err = io.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return
}

// GetMD5Hash return MD5 hash from text
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
