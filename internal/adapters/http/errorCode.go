package http

import (
	"context"
	"errors"
	"github.com/azargarov/rsvpck/internal/domain"
	"net"
	"strings"
	"syscall"
)

type httpErrorInfo struct {
	Status    domain.Status
	ErrorCode domain.ErrorCode
}

func classifyHTTPError(err, contextErr error) httpErrorInfo {
	if contextErr != nil {
		if errors.Is(contextErr, context.DeadlineExceeded) ||
			errors.Is(contextErr, context.Canceled) {
			return httpErrorInfo{
				Status:    domain.StatusTimeout,
				ErrorCode: domain.ErrorCodeTCPTimedOut,
			}
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return httpErrorInfo{
				Status:    domain.StatusTimeout,
				ErrorCode: domain.ErrorCodeTCPTimedOut,
			}
		}
	}

	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*syscall.Errno); ok {
			switch *syscallErr {
			case syscall.ECONNREFUSED:
				return httpErrorInfo{
					Status:    domain.StatusConnectionRefused,
					ErrorCode: domain.ErrorCodeConnectionRefused,
				}
			case syscall.ENETUNREACH, syscall.EHOSTUNREACH:
				return httpErrorInfo{
					Status:    domain.StatusTimeout,
					ErrorCode: domain.ErrorCodeTCPTimedOut,
				}
			}
		}
	}

	if strings.Contains(err.Error(), "tls:") || strings.Contains(err.Error(), "handshake") {
		if strings.Contains(err.Error(), "timeout") {
			return httpErrorInfo{
				Status:    domain.StatusTimeout,
				ErrorCode: domain.ErrorCodeTCPTimedOut,
			}
		}
		return httpErrorInfo{
			Status:    domain.StatusConnectionRefused,
			ErrorCode: domain.ErrorCodeConnectionRefused,
		}
	}

	if strings.Contains(err.Error(), "proxy") &&
		(strings.Contains(err.Error(), "auth") || strings.Contains(err.Error(), "authentication")) {
		return httpErrorInfo{
			Status:    domain.StatusProxyAuth,
			ErrorCode: domain.ErrorCodeProxyAuthRequired,
		}
	}

	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "connection refused") {
		return httpErrorInfo{
			Status:    domain.StatusConnectionRefused,
			ErrorCode: domain.ErrorCodeConnectionRefused,
		}
	}
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "context deadline exceeded") {
		return httpErrorInfo{
			Status:    domain.StatusTimeout,
			ErrorCode: domain.ErrorCodeTCPTimedOut,
		}
	}
	if strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "dns") {
		return httpErrorInfo{
			Status:    domain.StatusDNSFailure,
			ErrorCode: domain.ErrorCodeDNSUnresolvable,
		}
	}

	// Fallback: generic HTTP error
	return httpErrorInfo{
		Status:    domain.StatusHTTPError,
		ErrorCode: domain.ErrorCodeHTTPBadStatus,
	}
}
