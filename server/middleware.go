// This package creates and run a rest API server, using Gin framework
package server

import (
	"bytes"
	"encoding/json"
	"path"
	"strconv"

	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger send requests and response to database, and generate checksum if needed
func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Generate unique ID to make link between request and its associated response (stored in a different table)
		corID, _ := util.GenerateUID()

		body, err := util.ExtractBody(c.Request)
		if err != nil {
			logrus.Error("Logger coudn't send request data: ", err)
			return
		}

		//Create a checksum for the current request, only if it's the main endpoint (/fizzbuzz) and data is valid.
		var data model.Input
		checksum := ""
		if c.Request.RequestURI == URLPrefixVersion+FizzBaseURI && json.Unmarshal(body, &data) == nil && s.container.Validator.Struct(data) == nil {
			checksum = util.GetMD5Hash(data.String())
		}

		// Create copy to be used inside the goroutine - See Gin documentation : https://gin-gonic.com/docs/examples/goroutines-inside-a-middleware/
		cCp := c.Copy()

		go func() {
			if err := s.container.Repo.LogToDB("request", string(body), Scheme+"://"+path.Join(cCp.Request.Host, cCp.Request.RequestURI), corID, checksum, ""); err != nil {
				logrus.Error(err)
			}
		}()

		// Intercept Writer in order to get response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		status := blw.Status()
		respBody := blw.body.String()
		go func() {
			if err := s.container.Repo.LogToDB("response", respBody, "", corID, checksum, strconv.Itoa(status)); err != nil {
				logrus.Error(err)
			}
		}()
	}
}
