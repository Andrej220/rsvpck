package dns

import (    
	"context"
    "net"
    "net/netip"
    "time"
	"strings"
	"errors"
	"github.com/azargarov/rsvpck/internal/domain"
)

type Checker struct{}

var _ domain.DNSChecker = (*Checker)(nil)


func (r Checker) CheckWithContext(ctx context.Context, ep domain.Endpoint) domain.Probe {
	start := time.Now()
	_, err := net.DefaultResolver.LookupHost(ctx, ep.Target)
	latencyMs := time.Since(start).Seconds() * 1000

	if err != nil {
		return domain.NewFailedProbe(
			ep,
			mapDNSError(err, ctx.Err()),
			err,
		)
	}
	return domain.NewSuccessfulProbe(ep, latencyMs)
}


func (r *Checker)LookupHost(ctx context.Context, host string, timeout time.Duration) (ips []netip.Addr, err error){
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    strIPs, err := net.DefaultResolver.LookupHost(ctx, host)
    if err != nil {
        return nil, err
    }

    addrs := make([]netip.Addr, 0, len(strIPs))
    for _, ip := range strIPs {
        if addr, perr := netip.ParseAddr(ip); perr == nil {
            addrs = append(addrs, addr)
        }
    }

    return addrs, nil
}

func mapDNSError(err, contextErr error) domain.Status {
	// 1. First, check if the operation was cancelled or timed out via context
	if contextErr != nil {
		if errors.Is(contextErr, context.DeadlineExceeded) ||
			errors.Is(contextErr, context.Canceled) {
			return domain.StatusTimeout
		}
	}

	// 2. Handle *net.DNSError (most common DNS error type)
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		// If the DNS server didn't respond (timeout, network unreachable)
		if dnsErr.Timeout() {
			return domain.StatusTimeout
		}
		// If the host does not exist (NXDOMAIN)
		if dnsErr.IsNotFound {
			return domain.StatusDNSFailure // or a more specific StatusHostNotFound if you have it
		}
		// If no DNS server is configured or reachable
		if dnsErr.Err == "no such host" || dnsErr.Err == "server misbehaving" {
			return domain.StatusDNSFailure
		}
	}

	// 3. Fallback: check error message (less reliable, but safe)
	errStr := err.Error()
	if containsAny(errStr,
		"no such host",
		"server misbehaving",
		"cannot unmarshal DNS",
		"host not found",
		"NXDOMAIN",
	) {
		return domain.StatusDNSFailure
	}

	if containsAny(errStr,
		"timeout",
		"i/o timeout",
		"network is unreachable",
		"connection refused",
	) {
		return domain.StatusTimeout
	}

	// 4. Unknown error → treat as invalid or generic failure
	return domain.StatusInvalid
}

func containsAny(s string, substrs ...string) bool {
	sLower := strings.ToLower(s)
	for _, substr := range substrs {
		if substr != "" && strings.Contains(sLower, strings.ToLower(substr)) {
			return true
		}
	}
	return false
}