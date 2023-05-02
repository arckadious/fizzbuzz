// This package contains JSON response templates functions
package response

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	StatusError   = "error"
	StatusSuccess = "success"

	ErrorInvalidData         = "INVALID_DATA_ERR"
	ErrorInvalidField        = "INVALID_FIELD_ERR"
	ErrorInternalServerError = "INTERNAL_SERVER_ERR"
	ErrorRouteNotFound       = "ROUTE_NOT_FOUND_ERR"
	ErrorMethodNotAllowed    = "METHOD_NOT_ALLOWED_ERR"
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
	Field   string `json:"field"`
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

// SetStatusCode set api response status code
func (res *ApiResponse) SetStatusCode(statusCode int) *ApiResponse {
	res.StatusCode = statusCode
	return res
}

// SetStatus set api response content status
func (res *ApiResponse) SetStatus(status string) *ApiResponse {
	res.Status = status
	return res
}

// Setdata set api response content data
func (res *ApiResponse) SetData(data interface{}) *ApiResponse {
	res.Data = data
	return res
}

// SetErrorMessages set api response meta data
func (res *ApiResponse) SetErrorMessages(messages []ApiError) *ApiResponse {
	res.Messages = messages
	return res
}

// SetCustomErrorResponse set custom error response to be returned
func (res *ApiResponse) SetCustomErrorResponse(httpCode int, status string, messages []ApiError) *ApiResponse {
	res.SetStatusCode(httpCode)
	res.SetStatus(status)
	res.SetErrorMessages(messages)
	return res
}

// SetBadRequestResponse
func (res *ApiResponse) SetBadRequestResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusBadRequest, StatusError, messages)
	return res
}

// SetForbiddenResponse
func (res *ApiResponse) SetForbiddenResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusForbidden, StatusError, messages)
	return res
}

// SetNotFoundResponse
func (res *ApiResponse) SetNotFoundResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusNotFound, StatusError, messages)
	return res
}

// SetInternalServerErrorResponse
func (res *ApiResponse) SetInternalServerErrorResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusInternalServerError, StatusError, messages)
	return res
}

// AddErrorMessage add new error message to error list
func (res *ApiResponse) AddErrorMessage(msg ApiError) {
	res.Messages = append(res.Messages, msg)
}

// CreateErrorMessage returns new error struct from code and message
func (res *ApiResponse) CreateErrorMessage(code string, field string, message string) ApiError {
	return ApiError{code, field, message}
}

// GetErrorMessagesSLice returns an array of code with its description
func (res *ApiResponse) GetErrorMessageSlice(code string, field string, message string) []ApiError {
	var erros []ApiError
	erros = append(erros, res.CreateErrorMessage(code, field, message))
	return erros
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
