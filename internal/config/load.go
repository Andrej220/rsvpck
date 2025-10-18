package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/azargarov/rsvpck/internal/domain"
)

type FileSpec struct {
	ProxyURL        string         `json:"proxyURL"`
	VPNEndpoints    []EndpointSpec `json:"vpnEndpoints"`
	DirectEndpoints []EndpointSpec `json:"directEndpoints"`
	ProxyEndpoints  []EndpointSpec `json:"proxyEndpoints"`
}

type EndpointSpec struct {
	Target   string `json:"target"`   // e.g. "google.com:443" or "https://..."
	Type     string `json:"type"`     // "public" | "vpn"
	Kind     string `json:"kind"`     // "icmp" | "dns" | "tcp" | "http"
	Note     string `json:"note"`     // human description
	UseProxy bool   `json:"useProxy"` // for HTTP kind
}

// LoadFromFile builds a NetTestConfig from a JSON file.
// If path is empty, returns an error so caller can fallback to Default.
func LoadFromFile(path string) (domain.NetTestConfig, error) {
	if path == "" {
		return domain.NetTestConfig{}, errors.New("empty path")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return domain.NetTestConfig{}, err
	}

	var spec FileSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return domain.NetTestConfig{}, err
	}

	toEP := func(s EndpointSpec) domain.Endpoint {
		etype := domain.EndpointTypePublic
		if s.Type == "vpn" {
			etype = domain.EndpointTypeVPN
		}
		switch s.Kind {
		case "icmp":
			return domain.MustNewICMPEndpoint(s.Target, etype, s.Note)
		case "dns":
			return domain.MustNewDNSEndpoint(s.Target, etype, s.Note)
		case "tcp":
			return domain.MustNewTCPEndpoint(s.Target, etype, s.Note)
		case "http":
			return domain.MustNewHTTPEndpoint(s.Target, etype, s.UseProxy, spec.ProxyURL, s.Note)
		default:
			panic("unknown endpoint kind: " + s.Kind)
		}
	}

	var vpn, direct, proxy []domain.Endpoint
	for _, e := range spec.VPNEndpoints {
		vpn = append(vpn, toEP(e))
	}
	for _, e := range spec.DirectEndpoints {
		direct = append(direct, toEP(e))
	}
	for _, e := range spec.ProxyEndpoints {
		proxy = append(proxy, toEP(e))
	}

	return domain.NewNetTestConfig(vpn, direct, proxy, spec.ProxyURL)
}
