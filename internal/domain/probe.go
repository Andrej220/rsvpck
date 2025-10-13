package domain

import (
	"fmt"
	"time"
)

type Probe struct {
	Endpoint  Endpoint
	Status    Status
	LatencyMs float64
	Error     string
	Timestamp time.Time
}

func (p Probe) IsSuccessful() bool {
	return p.Status == StatusPass
}

func (p Probe) IsSkipped() bool {
	return p.Status == StatusSkipped
}

func (p Probe) IsVPNProbe() bool {
	return p.Endpoint.IsVPN()
}

func (p Probe) IsDNSProbe() bool {
	return p.Endpoint.IsDNS()
}

func (p Probe) String() string {
	var str string
	if !p.IsSuccessful() {
		str = fmt.Sprintf("Endpoint: %s, Status: %s, Error: %s, ",
			p.Endpoint.String(), p.Status.String(), p.Error)
		return str

	}
	str = fmt.Sprintf("Endpoint: %v, Status: %s, Latency: %0.2f, ",
		p.Endpoint, p.Status.String(), p.LatencyMs)

	return str
}

func NewSuccessfulProbe(endpoint Endpoint, latencyMs float64) Probe {
	return Probe{Status: StatusPass, Endpoint: endpoint, LatencyMs: latencyMs}
}
func NewFailedProbe(endpoint Endpoint, status Status, err error) Probe {
	return Probe{Status: StatusFail, Endpoint: endpoint, Error: err.Error()}
}
