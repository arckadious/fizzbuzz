package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusError   = "error"
	StatusSuccess = "success"
	StatusWarning = "warning"

	ErrorInvalidData         = "INVALID_DATA_ERR"
	ErrorInvalidField        = "INVALID_FIELD_ERR"
	ErrorInternalServerError = "INTERNAL_SERVER_ERR"
	ErrorRouteNotFound       = "ROUTE_NOT_FOUND_ERR"
	ErrorMethodNotAllowed    = "METHOD_NOT_ALLOWED_ERR"
)

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

// Constructor of the new error
func NewApiError(code string, message string, field string) *ApiError {
	return &ApiError{
		Code:    code,
		Field:   field,
		Message: message,
	}
}

// Redefine the error interface function
func (c *ApiError) Error() string {
	return c.Message
}

// New create a new Response
func NewApiResponse(statusCode int, status string, errorMessages []ApiError, data interface{}) *ApiResponse {
	return &ApiResponse{
		StatusCode: statusCode,
		Status:     status,
		Messages:   errorMessages,
		Data:       data,
	}
}

// Set api response status code
func (res *ApiResponse) SetStatusCode(statusCode int) *ApiResponse {
	res.StatusCode = statusCode
	return res
}

// Set api response content status
func (res *ApiResponse) SetStatus(status string) *ApiResponse {
	res.Status = status
	return res
}

// Set api response content data
func (res *ApiResponse) SetData(data interface{}) *ApiResponse {
	res.Data = data
	return res
}

// Set api response meta data
func (res *ApiResponse) SetErrorMessages(messages []ApiError) *ApiResponse {
	res.Messages = messages
	return res
}

// Set custom error response to be returned
func (res *ApiResponse) SetCustomErrorResponse(httpCode int, status string, messages []ApiError) *ApiResponse {
	res.SetStatusCode(httpCode)
	res.SetStatus(status)
	res.SetErrorMessages(messages)
	return res
}

// Set automatique bad request response
func (res *ApiResponse) SetBadRequestResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusBadRequest, StatusError, messages)
	return res
}

// Set automatique forbidden response
func (res *ApiResponse) SetForbiddenResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusForbidden, StatusError, messages)
	return res
}

// Set automatique bad request response
func (res *ApiResponse) SetNotFoundResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusNotFound, StatusError, messages)
	return res
}

// Set automatique bad request response
func (res *ApiResponse) SetInternalServerErrorResponse(messages []ApiError) *ApiResponse {
	res.SetCustomErrorResponse(http.StatusInternalServerError, StatusError, messages)
	return res
}

// Add new error message to error list
func (res *ApiResponse) AddErrorMessage(msg ApiError) {
	res.Messages = append(res.Messages, msg)
}

// Returns new error struct from code and message
func (res *ApiResponse) CreateErrorMessage(code string, field string, message string) ApiError {
	return ApiError{code, field, message}
}

func (res *ApiResponse) GetErrorMessageSlice(code string, field string, message string) []ApiError {
	var erros []ApiError
	erros = append(erros, res.CreateErrorMessage(code, field, message))
	return erros
}

// WriteJson output write json response
func (res *ApiResponse) WriteJsonResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	resp, _ := json.Marshal(res)
	w.Write(resp)
}
