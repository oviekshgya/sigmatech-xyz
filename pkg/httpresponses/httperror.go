package httpresponses

import (
	"net/http"
)

type HttpError struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
	Error   string `json:"error"`
}

func NewBadRequestError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusBadRequest,
		Error:   "bad_request",
	}
}
func NewUnauthorizedError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusUnauthorized,
		Error:   "UnAuthorized",
	}
}
func NewNotFoundError(message string) *HttpError {
	//LogErrors(message, strconv.Itoa(http.StatusNotFound))
	return &HttpError{
		Message: message,
		Code:    http.StatusNotFound,
		Error:   "not_found",
	}
}

func NewSuccess(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusOK,
		Error:   "Sukses",
	}
}

func NewInternalServerError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}
