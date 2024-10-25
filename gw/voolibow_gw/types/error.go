// Package types defines application-specific error types and conversion functions for gRPC status.

package types

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error represents an application-specific error with a message and code.
type Error struct {
	Message string // Error message
	Code    int    // Error code
}

// NewError creates a new custom error with the given code and message.
func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// NewInternalError creates a new internal server error with the given message.
func NewInternalError(message string) *Error {
	return &Error{Code: 500, Message: message}
}

// NewNotFoundError creates a new not found error with the given message.
func NewNotFoundError(message string) *Error {
	return &Error{Code: 404, Message: message}
}

// NewPermissionDeniedError creates a new permission denied error with the given message.
func NewPermissionDeniedError(message string) *Error {
	return &Error{Code: 403, Message: message}
}

// NewBadRequestError creates a new bad request error with the given message.
func NewBadRequestError(message string) *Error {
	return &Error{Code: 400, Message: message}
}

// NewUnauthorizedError creates an unauthorized request error with the given message.
func NewUnauthorizedError(message string) *Error {
	return &Error{Code: 401, Message: message}
}

// ErrorToGRPCStatus converts the custom error to a gRPC status error based on the error code.
func (c *Error) ErrorToHttpError() error {
	switch c.Code {
	case 500:
		return status.Error(codes.Internal, c.Message)
	case 404:
		return status.Error(codes.NotFound, c.Message)
	case 403:
		return status.Error(codes.PermissionDenied, c.Message)
	case 400:
		return status.Error(codes.Aborted, c.Message)
	default:
		return status.Error(codes.Aborted, c.Message)
	}
}

func (c *Error) ErrorToWebSocketState() string {
	switch c.Code {
	case 500:
		return "internal"
	case 404:
		return "not_found"
	case 403:
		return "permission_denied"
	case 400:
		return "bad_request"
	default:
		return "bad_request"
	}
}

// ExtractGrpcError extracts a custom error from a gRPC status error.
func ExtractGrpcError(err error) *Error {
	if err == nil {
		return NewError(500, "unknown error")
	}

	fmt.Println(err)
	st, ok := status.FromError(err)
	if !ok {
		return NewError(500, "unknown error")
	}

	grpcCode := st.Code()

	switch grpcCode {
	case codes.Aborted:
		return NewBadRequestError(st.Message())
	case codes.Internal:
		return NewInternalError(st.Message())
	case codes.NotFound:
		return NewNotFoundError(st.Message())
	case codes.PermissionDenied:
		return NewPermissionDeniedError(st.Message())
	default:
		return NewError(500, "unknown error")
	}
}
