// This package contains JSON response template functions
package response

import (
	"encoding/json"
	"net/http"

	cst "github.com/arckadious/fizzbuzz/constant"
	"github.com/sirupsen/logrus"
)

// ApiResponse class
type ApiResponse struct {
	StatusCode int         `json:"-"`
	Status     string      `json:"status"`
	Messages   []ApiError  `json:"messages"`
	Data       interface{} `json:"data"`
}

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// New constructor ApiResponse : create a new response template
func New(statusCode int, status string, errorMessages []ApiError, data interface{}) *ApiResponse {
	return &ApiResponse{
		StatusCode: statusCode,
		Status:     status,
		Messages:   errorMessages,
		Data:       data,
	}
}

// Setdata set api response content data
func (res *ApiResponse) SetData(data interface{}) *ApiResponse {
	res.Data = data
	return res
}

// SetErrorResponse set custom error response to be returned
func (res *ApiResponse) SetErrorResponse(statusCode int, messages []ApiError) *ApiResponse {
	res.StatusCode = statusCode
	res.Status = cst.StatusError
	res.Messages = messages
	return res
}

// WriteJSONResponse writes json output to json response
func (res *ApiResponse) WriteJSONResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	resp, err := json.Marshal(res)
	if err != nil {
		logrus.Error(err)
	}
	w.Write(resp)
}
