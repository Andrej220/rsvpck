package domain

import "time"

type ConnectivityResult struct {
	Mode        ConnectivityMode
	IsConnected bool
	Probes      []Probe
	Timestamp   time.Time
	Summary     string
}

func NewConnectivityResult(mode ConnectivityMode, probes []Probe) ConnectivityResult {
	isConnected := mode != ModeNone
	summary := buildSummary(mode, probes)
	return ConnectivityResult{
		Mode:        mode,
		IsConnected: isConnected,
		Probes:      probes,
		Timestamp:   time.Now(),
		Summary:     summary,
	}
}

func (r ConnectivityResult) SuccessfulProbes() []Probe {
	var success []Probe
	for _, p := range r.Probes {
		if p.IsSuccessful() {
			success = append(success, p)
		}
	}
	return success
}

func (r ConnectivityResult) FailedProbes() []Probe {
	var failed []Probe
	for _, p := range r.Probes {
		if !p.IsSuccessful() {
			failed = append(failed, p)
		}
	}
	return failed
}
