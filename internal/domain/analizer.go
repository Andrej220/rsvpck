package domain

func AnalyzeConnectivity(probes []Probe, config NetTestConfig) ConnectivityResult {
	var (
		vpnOK      bool
		directOK   bool
		proxyOK    bool
		dnsOK bool
	)

	for _, p := range probes {
		if p.IsDNSProbe() && p.IsSuccessful() {
			dnsOK = true
			break
		}
	}

	for _, p := range probes {
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

	var mode ConnectivityMode
	switch {
	case directOK && dnsOK:
		mode = ModeDirect
	case proxyOK:
		mode = ModeViaProxy
	case vpnOK:
		mode = ModeViaVPN
	default:
		mode = ModeNone
	}

	return NewConnectivityResult(mode, probes)
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
