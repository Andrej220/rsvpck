// internal/domain/errors.go
package domain

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	ErrorCodeDNSUnresolvable ErrorCode = iota
	ErrorCodeTCPTimedOut
	ErrorCodeHTTPBadStatus
	ErrorCodeConnectionRefused
	ErrorCodeProxyAuthRequired
	ErrorCodeInvalidConfig
	ErrorCodeICMPFailed
	ErrorCodeHTTPClientError
)

func (ec ErrorCode) Error() string {
	switch ec {
	case ErrorCodeDNSUnresolvable:
		return "DNS resolution failed"
	case ErrorCodeTCPTimedOut:
		return "TCP connection timed out"
	case ErrorCodeHTTPBadStatus:
		return "HTTP request returned non-success status"
	case ErrorCodeConnectionRefused:
		return "connection refused"
	case ErrorCodeProxyAuthRequired:
		return "proxy authentication required"
	case ErrorCodeInvalidConfig:
		return "invalid configuration"
	case ErrorCodeICMPFailed:
		return "ICMP ping failed"
	default:
		return fmt.Sprintf("unknown error code: %d", ec)
	}
}

func (ec ErrorCode) Is(target error) bool {
	if targetCode, ok := target.(ErrorCode); ok {
		return ec == targetCode
	}
	return false
}

func Errorf(code ErrorCode, format string, args ...any) error {
	return &wrappedError{
		code:    code,
		inner: fmt.Errorf(format, args...),
	}
}

type wrappedError struct {
	code    ErrorCode
	inner   error
}

func (e *wrappedError) Error() string {
	return e.inner.Error()
}

func (e *wrappedError) Unwrap() error {
	return e.inner
}

func IsErrorCode(err error, code ErrorCode) bool {
	return errors.Is(err, code)
}

func (e *wrappedError) Is(target error) bool {
	if targetCode, ok := target.(ErrorCode); ok {
		return e.code == targetCode
	}
	return errors.Is(e.inner, target)
}

func ErrInvalidConfig(reason string) error {
	return Errorf(ErrorCodeInvalidConfig, "invalid config: %s", reason)
}