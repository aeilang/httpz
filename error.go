// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: © 2015 LabStack LLC and Echo contributors
// copied from echo, source: https://github.com/labstack/echo

package httpz

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// ErrHandlerFunc defines the function signature for centralized error handling.
type ErrHandlerFunc func(err error, w http.ResponseWriter)

// DefaultErrHandlerFunc is the default centralized error handling function.
// It only triggers an error response for *HTTPError.
func DefaultErrHandlerFunc(err error, w http.ResponseWriter) {
	if he, ok := err.(*HTTPError); ok {
		rw := NewHelperRW(w)
		rw.JSON(he.StatusCode, Map{"msg": he.Msg})
	} else {
		slog.Error(err.Error())
	}
}

// HTTPError represents a custom error type inspired by Echo.
type HTTPError struct {
	StatusCode int    // HTTP status code
	Msg        string // Error message
	Internal   error  // Internal error
}

// NewHTTPError creates a new HTTPError with the given status code and message.
func NewHTTPError(statusCode int, msg string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

// SetInternal sets the internal error for the HTTPError.
func (e *HTTPError) SetInternal(err error) *HTTPError {
	e.Internal = err
	return e
}

// Error returns the error message for the HTTPError.
func (e *HTTPError) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", e.StatusCode, e.Msg)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", e.StatusCode, e.Msg, e.Internal)
}

// Unwrap returns the internal error of the HTTPError.
func (e *HTTPError) Unwrap() error {
	return e.Internal
}

// helper returns the status code and its corresponding text.
func helper(code int) (int, string) {
	return code, http.StatusText(code)
}

// Predefined HTTP errors
var (
	ErrBadRequest                    = NewHTTPError(helper(http.StatusBadRequest))                    // HTTP 400 Bad Request
	ErrUnauthorized                  = NewHTTPError(helper(http.StatusUnauthorized))                  // HTTP 401 Unauthorized
	ErrPaymentRequired               = NewHTTPError(helper(http.StatusPaymentRequired))               // HTTP 402 Payment Required
	ErrForbidden                     = NewHTTPError(helper(http.StatusForbidden))                     // HTTP 403 Forbidden
	ErrNotFound                      = NewHTTPError(helper(http.StatusNotFound))                      // HTTP 404 Not Found
	ErrMethodNotAllowed              = NewHTTPError(helper(http.StatusMethodNotAllowed))              // HTTP 405 Method Not Allowed
	ErrNotAcceptable                 = NewHTTPError(helper(http.StatusNotAcceptable))                 // HTTP 406 Not Acceptable
	ErrProxyAuthRequired             = NewHTTPError(helper(http.StatusProxyAuthRequired))             // HTTP 407 Proxy AuthRequired
	ErrRequestTimeout                = NewHTTPError(helper(http.StatusRequestTimeout))                // HTTP 408 Request Timeout
	ErrConflict                      = NewHTTPError(helper(http.StatusConflict))                      // HTTP 409 Conflict
	ErrGone                          = NewHTTPError(helper(http.StatusGone))                          // HTTP 410 Gone
	ErrLengthRequired                = NewHTTPError(helper(http.StatusLengthRequired))                // HTTP 411 Length Required
	ErrPreconditionFailed            = NewHTTPError(helper(http.StatusPreconditionFailed))            // HTTP 412 Precondition Failed
	ErrStatusRequestEntityTooLarge   = NewHTTPError(helper(http.StatusRequestEntityTooLarge))         // HTTP 413 Payload Too Large
	ErrRequestURITooLong             = NewHTTPError(helper(http.StatusRequestURITooLong))             // HTTP 414 URI Too Long
	ErrUnsupportedMediaType          = NewHTTPError(helper(http.StatusUnsupportedMediaType))          // HTTP 415 Unsupported Media Type
	ErrRequestedRangeNotSatisfiable  = NewHTTPError(helper(http.StatusRequestedRangeNotSatisfiable))  // HTTP 416 Range Not Satisfiable
	ErrExpectationFailed             = NewHTTPError(helper(http.StatusExpectationFailed))             // HTTP 417 Expectation Failed
	ErrTeapot                        = NewHTTPError(helper(http.StatusTeapot))                        // HTTP 418 I'm a teapot
	ErrMisdirectedRequest            = NewHTTPError(helper(http.StatusMisdirectedRequest))            // HTTP 421 Misdirected Request
	ErrUnprocessableEntity           = NewHTTPError(helper(http.StatusUnprocessableEntity))           // HTTP 422 Unprocessable Entity
	ErrLocked                        = NewHTTPError(helper(http.StatusLocked))                        // HTTP 423 Locked
	ErrFailedDependency              = NewHTTPError(helper(http.StatusFailedDependency))              // HTTP 424 Failed Dependency
	ErrTooEarly                      = NewHTTPError(helper(http.StatusTooEarly))                      // HTTP 425 Too Early
	ErrUpgradeRequired               = NewHTTPError(helper(http.StatusUpgradeRequired))               // HTTP 426 Upgrade Required
	ErrPreconditionRequired          = NewHTTPError(helper(http.StatusPreconditionRequired))          // HTTP 428 Precondition Required
	ErrTooManyRequests               = NewHTTPError(helper(http.StatusTooManyRequests))               // HTTP 429 Too Many Requests
	ErrRequestHeaderFieldsTooLarge   = NewHTTPError(helper(http.StatusRequestHeaderFieldsTooLarge))   // HTTP 431 Request Header Fields Too Large
	ErrUnavailableForLegalReasons    = NewHTTPError(helper(http.StatusUnavailableForLegalReasons))    // HTTP 451 Unavailable For Legal Reasons
	ErrInternalServerError           = NewHTTPError(helper(http.StatusInternalServerError))           // HTTP 500 Internal Server Error
	ErrNotImplemented                = NewHTTPError(helper(http.StatusNotImplemented))                // HTTP 501 Not Implemented
	ErrBadGateway                    = NewHTTPError(helper(http.StatusBadGateway))                    // HTTP 502 Bad Gateway
	ErrServiceUnavailable            = NewHTTPError(helper(http.StatusServiceUnavailable))            // HTTP 503 Service Unavailable
	ErrGatewayTimeout                = NewHTTPError(helper(http.StatusGatewayTimeout))                // HTTP 504 Gateway Timeout
	ErrHTTPVersionNotSupported       = NewHTTPError(helper(http.StatusHTTPVersionNotSupported))       // HTTP 505 HTTP Version Not Supported
	ErrVariantAlsoNegotiates         = NewHTTPError(helper(http.StatusVariantAlsoNegotiates))         // HTTP 506 Variant Also Negotiates
	ErrInsufficientStorage           = NewHTTPError(helper(http.StatusInsufficientStorage))           // HTTP 507 Insufficient Storage
	ErrLoopDetected                  = NewHTTPError(helper(http.StatusLoopDetected))                  // HTTP 508 Loop Detected
	ErrNotExtended                   = NewHTTPError(helper(http.StatusNotExtended))                   // HTTP 510 Not Extended
	ErrNetworkAuthenticationRequired = NewHTTPError(helper(http.StatusNetworkAuthenticationRequired)) // HTTP 511 Network Authentication Required

	ErrValidatorNotRegistered = errors.New("validator not registered")
	ErrRendererNotRegistered  = errors.New("renderer not registered")
	ErrInvalidRedirectCode    = errors.New("invalid redirect status code")
	ErrCookieNotFound         = errors.New("cookie not found")
	ErrInvalidCertOrKeyType   = errors.New("invalid cert or key type, must be string or []byte")
	ErrInvalidListenerNetwork = errors.New("invalid listener network")
)
