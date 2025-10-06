package domain


func AnalyzeConnectivity(probes []Probe, config NetTestConfig) ConnectivityResult {
	var (
		vpnSuccess     bool
		directSuccess  bool
		proxySuccess   bool
		dnsSuccess     bool
	)

	for _, p := range probes {
		if p.IsSuccessful() {
			if p.Endpoint.IsVPN() {
				vpnSuccess = true
			} else if p.Endpoint.IsPublic() {
				if p.Endpoint.TargetType == TargetTypeTCP {
					directSuccess = true
				} else if p.Endpoint.TargetType == TargetTypeHTTP && p.Endpoint.MustUseProxy() {
					proxySuccess = true
				}
			}
		}
	}

	for _, p := range probes {
		if p.IsDNSProbe() && p.IsSuccessful() {
			dnsSuccess = true
			break
		}
	}

	var mode ConnectivityMode
	if vpnSuccess {
		mode = ModeViaVPN
	} else if directSuccess && (dnsSuccess || !config.HasDNSCheck()) {
		mode = ModeDirect
	} else if proxySuccess {
		mode = ModeViaProxy
	} else {
		mode = ModeNone
	}

	return NewConnectivityResult(mode, probes)
}

func buildSummary(mode ConnectivityMode, probes []Probe) string {
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