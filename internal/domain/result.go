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
	summary := buildSummary(mode)
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

func (r *ConnectivityResult) DetermineMode() {
	var (
		vpnOK, directOK, proxyOK, dnsOK bool
	)

	for _, p := range r.Probes {
		if p.IsDNSProbe() && p.IsSuccessful() {
			dnsOK = true
			break
		}
	}

	for _, p := range r.Probes {
		if !p.IsSuccessful() {
			continue
		}

		switch {
		case p.Endpoint.IsVPN():
			vpnOK = true

		case !p.Endpoint.IsPublic():
			continue

		case p.Endpoint.IsDirectType():
			directOK = true

		case p.Endpoint.IsProxyType():
			proxyOK = true
		}
	}

	switch {
	case directOK && dnsOK:
		r.Mode = ModeDirect
	case proxyOK:
		r.Mode = ModeViaProxy
	case vpnOK:
		r.Mode = ModeViaVPN
	default:
		r.Mode = ModeNone
	}

	r.IsConnected = r.Mode.IsConnected()
	r.Summary = buildSummary(r.Mode)
}

func buildSummary(mode ConnectivityMode) string {
	switch mode {
	case ModeViaVPN:
		return "Connected via VPN."
	case ModeDirect:
		return "Direct internet."
	case ModeViaProxy:
		return "Internet via proxy"
	default:
		return "No connection"
	}
}