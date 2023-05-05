// This package creates and run a rest API server, using Gin framework
package server

import (
	"net/http"

	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/arckadious/fizzbuzz/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// recoveryHandler handles unexpected 'panic' call during process, and return a custom 500 internal server error.
func (s *Server) recoveryHandler(c *gin.Context, err interface{}) {
	logrus.Error(err)
	response.New(
		http.StatusInternalServerError,
		cst.StatusError,
		[]response.ApiError{
			{
				Code:    cst.ErrorInternalServerError,
				Message: "Internal Server Error, oups !",
			},
		},
		nil,
	).WriteJSONResponse(c.Writer)
}

// notFoundHandler handles API server response when endpoint couldn't be found
func (s *Server) notFoundHandler(c *gin.Context) {
	response.New(
		http.StatusNotFound,
		cst.StatusError,
		[]response.ApiError{
			{
				Code:    cst.ErrorRouteNotFound,
				Message: "This is not what you are looking for.",
			},
		},
		nil,
	).WriteJSONResponse(c.Writer)
}

// methodNotAllowedHandler handles API server response when endpoint method is not allowed
func (s *Server) methodNotAllowedHandler(c *gin.Context) {
	response.New(
		http.StatusMethodNotAllowed,
		cst.StatusError,
		[]response.ApiError{
			{
				Code:    cst.ErrorMethodNotAllowed,
				Message: "Method is not allowed, boy.",
			},
		},
		nil,
	).WriteJSONResponse(c.Writer)
}
