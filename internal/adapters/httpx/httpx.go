package httpx

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/azargarov/rsvpck/internal/domain"
	"net"
	"strings"
	"time"
)

// TLSChecker implements domain.TLSCertificateFetcher.
//
// Compile-time interface assertion:
//var _ domain.TLSCertificateFetcher = (*TLSChecker)(nil)

func GetCertificates(ctx context.Context, addr, serverName string) ([]domain.TLSCertificate, error) {

	var timeout time.Duration = 1 * time.Second

	if addr == "" {
		return nil, errors.New("addr is required (host:port)")
	}

	if serverName == "" {
		serverName = hostPart(addr)
	}

	dialer := &net.Dialer{}
	if deadline, ok := ctx.Deadline(); ok {
		dialer.Timeout = time.Until(deadline)
	} else if timeout > 0 {
		dialer.Timeout = timeout
	}

	rawConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rawConn.Close() }()

	cfg := &tls.Config{
		ServerName: serverName, // SNI + hostname verification
	}
	tlsConn := tls.Client(rawConn, cfg)

	// Context-aware handshake (Go 1.20+).
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return nil, err
	}

	// After handshake succeeds, the underlying rawConn is now owned by tlsConn.
	// Close tlsConn instead.
	defer func() { _ = tlsConn.Close() }()
	state := tlsConn.ConnectionState()
	now := time.Now()

	out := make([]domain.TLSCertificate, 0, len(state.PeerCertificates))
	for _, cert := range state.PeerCertificates {
		out = append(out, domain.TLSCertificate{
			Subject:   cert.Subject.String(),
			Issuer:    cert.Issuer.String(),
			NotBefore: cert.NotBefore,
			NotAfter:  cert.NotAfter,
			Valid:     !now.Before(cert.NotBefore) && !now.After(cert.NotAfter),
		})
	}
	return out, nil
}

func hostPart(addr string) string {
	if h, _, err := net.SplitHostPort(addr); err == nil {
		return h
	}
	// Fallback: strip after first colon
	if i := strings.IndexByte(addr, ':'); i >= 0 {
		return addr[:i]
	}
	return addr
}
