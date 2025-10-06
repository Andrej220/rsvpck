package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"
	"github.com/azargarov/rsvpck/internal/domain"
)

type Checker struct{}

var _ domain.HTTPChecker = (*Checker)(nil)

func (c Checker) CheckWithContext(ctx context.Context, ep domain.Endpoint) domain.Probe {
	return c.doRequest(ctx, ep, nil)
}

func (c Checker) CheckViaProxyWithContext(ctx context.Context, ep domain.Endpoint, proxyURL string) domain.Probe {
	proxyParsed, err := url.Parse(proxyURL)
	if err != nil {
		return domain.NewFailedProbe(
			ep,
			domain.StatusInvalid,
			errors.New("invalid proxy URL: "+err.Error()),
		)
	}
	return c.doRequest(ctx, ep, proxyParsed)
}

// TODO: test with proxy
// doRequest is the shared HTTP execution logic.
func (c Checker) doRequest(ctx context.Context, ep domain.Endpoint, proxyURL *url.URL) domain.Probe {

	//proxyURL = nil
	var transport http.RoundTripper = http.DefaultTransport
	if proxyURL != nil {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.Proxy = func(*http.Request) (*url.URL, error) {
			return proxyURL, nil
		}
		transport = t
	}
	
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // overall request timeout
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // don't follow redirects
		},
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", ep.Target, nil)
	if err != nil {
		return domain.NewFailedProbe(
			ep,
			domain.StatusInvalid,
			err,
		)
	}
	
	// Add a minimal user-agent to avoid 403s
	req.Header.Set("User-Agent", "rsvpck/0.2 (network tester)")
	
	start := time.Now()
	resp, err := client.Do(req)
	latencyMs := time.Since(start).Seconds() * 1000

	if err != nil {
		status := mapHTTPError(err, ctx.Err())
		return domain.NewFailedProbe(
			ep,
			status,
			err,
		)
	}
	defer resp.Body.Close()

	// TODO: Consider 2xx and 3xx (redirects) as success for connectivity
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return domain.NewSuccessfulProbe(
			ep,
			latencyMs,
		)
	}

	// HTTP error (4xx, 5xx)
	return domain.NewFailedProbe(
		ep,
		domain.StatusHTTPError,
		errors.New("HTTP "+resp.Status),
	)
}

func mapHTTPError(err, contextErr error) domain.Status {
	if contextErr != nil {
		if errors.Is(contextErr, context.DeadlineExceeded) ||
			errors.Is(contextErr, context.Canceled) {
			return domain.StatusTimeout
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return domain.StatusTimeout
		}
	}

	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*syscall.Errno); ok {
			switch *syscallErr {
			case syscall.ECONNREFUSED:
				return domain.StatusConnectionRefused
			case syscall.ENETUNREACH, syscall.EHOSTUNREACH:
				return domain.StatusTimeout
			}
		}
	}

	if strings.Contains(err.Error(), "tls:") || strings.Contains(err.Error(), "handshake") {
		if strings.Contains(err.Error(), "timeout") {
			return domain.StatusTimeout
		}
		return domain.StatusConnectionRefused
	}

	if strings.Contains(err.Error(), "proxy") &&
		(strings.Contains(err.Error(), "auth") || strings.Contains(err.Error(), "authentication")) {
		return domain.StatusProxyAuth
	}

	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "connection refused") {
		return domain.StatusConnectionRefused
	}
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "context deadline exceeded") {
		return domain.StatusTimeout
	}
	if strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "dns") {
		return domain.StatusDNSFailure
	}

	// Unknown error
	return domain.StatusInvalid
}