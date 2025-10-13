package domain

type ConnectivityMode int

const (
	ModeNone ConnectivityMode = iota
	ModeDirect
	ModeViaProxy
	ModeViaVPN
)

func (m ConnectivityMode) String() string {
	switch m {
	case ModeNone:
		return "none"
	case ModeDirect:
		return "direct"
	case ModeViaProxy:
		return "via_proxy"
	case ModeViaVPN:
		return "via_vpn"
	default:
		return "unknown"
	}
}

func (m ConnectivityMode) IsConnected() bool {
	return m == ModeDirect || m == ModeViaProxy || m == ModeViaVPN
}
