package domain

import (
	"context"
	"io"
)

type TCPChecker interface {
	CheckWithContext(ctx context.Context, ep Endpoint) Probe
}

type DNSChecker interface {
	CheckWithContext(ctx context.Context, ep Endpoint) Probe
}

type HTTPChecker interface {
	CheckWithContext(ctx context.Context, ep Endpoint) Probe                    
	CheckViaProxyWithContext(ctx context.Context, ep Endpoint, proxyURL string) Probe 
}

type ICMPChecker interface{
	CheckPingWithContext(ctx context.Context, ep Endpoint) Probe
}

type Renderer interface {
	Render(w io.Writer, result ConnectivityResult) error
}