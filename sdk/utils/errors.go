package utils

import "errors"

var (
	ErrInvalidConfiguration = errors.New("invalid client configuration")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrNotFound             = errors.New("resource not found")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrAPIError             = errors.New("API error")
)
