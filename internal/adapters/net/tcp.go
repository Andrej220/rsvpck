package netadapter

import (
    "context"
    "net"
    "time"
	"fmt"
    "github.com/azargarov/rsvpck/internal/domain"
)

type TCPDialer struct{}

var _ domain.TCPChecker = (*TCPDialer)(nil)

func (d *TCPDialer) CheckWithContext(ctx context.Context, ep domain.Endpoint) domain.Probe {
    addr := ep.Target //net.JoinHostPort(host, fmt.Sprint(port))
    start := time.Now()
    conn, err := net.DialTimeout("tcp", addr, 3*time.Second) 
    latencyMs := time.Since(start).Seconds() * 1000

    
    if err != nil {
        return domain.NewFailedProbe(
			ep,
			domain.StatusConnectionRefused,//mapPingError(err, ctx.Err(), output), //TODO: implement error mapping function
			err,
		)
    }
    conn.Close()

    return domain.NewSuccessfulProbe(ep, latencyMs)
}

func (d *TCPDialer) HTTPDo(ctx context.Context, req *domain.Request) (*domain.Response, time.Duration, error) {
    return nil, 0, fmt.Errorf("not implemented")
}

func (d *TCPDialer) ICMPPing(ctx context.Context, host string, count int) (time.Duration, error) {
    return 0, fmt.Errorf("not implemented")
}
