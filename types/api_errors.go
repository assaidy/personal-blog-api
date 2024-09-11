package types

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("api error: %d - %v", e.StatusCode, e.Msg)
}

func NewApiError(statusCode int, err error) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidJSONError() ApiError {
	return NewApiError(http.StatusBadRequest, fmt.Errorf("invalid JSON request data"))
}

func NotFoundError(err error) ApiError {
	return NewApiError(http.StatusNotFound, err)
}

func AlreadyExistsError(err error) ApiError {
	return NewApiError(http.StatusBadRequest, err)
}
