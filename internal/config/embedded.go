package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"github.com/azargarov/rsvpck/internal/domain"
)

//go:embed defaults/geconfig.json
var defaultsFS embed.FS

func LoadEmbedded() (domain.NetTestConfig, error) {
	b, err := defaultsFS.ReadFile("defaults/geconfig.json")
	if err != nil {
		return domain.NetTestConfig{}, err
	}
	return parseJSON(b)
}

func LoadFromFileOrEmbedded(path string) (domain.NetTestConfig, error) {
	if path != "" {
		if b, err := os.ReadFile(path); err == nil {
			return parseJSON(b)
		} else {
			return domain.NetTestConfig{}, fmt.Errorf("read %s: %w", path, err)
		}
	}
	return LoadEmbedded()
}

// parseJSON takes raw JSON bytes from an embedded or file config
// and converts them into a domain.NetTestConfig that your app uses.
func parseJSON(b []byte) (domain.NetTestConfig, error) {
	var spec FileSpec
	if err := json.Unmarshal(b, &spec); err != nil {
		return domain.NetTestConfig{}, fmt.Errorf("invalid JSON config: %w", err)
	}
	return specToDomain(spec)
}

// specToDomain converts a decoded FileSpec to a fully built domain.NetTestConfig.
func specToDomain(spec FileSpec) (domain.NetTestConfig, error) {
	toEndpoint := func(s EndpointSpec) (domain.Endpoint, error) {
		etype := domain.EndpointTypePublic
		if s.Type == "vpn" {
			etype = domain.EndpointTypeVPN
		}

		switch s.Kind {
		case "icmp":
			return domain.MustNewICMPEndpoint(s.Target, etype, s.Note), nil
		case "dns":
			return domain.MustNewDNSEndpoint(s.Target, etype, s.Note), nil
		case "tcp":
			return domain.MustNewTCPEndpoint(s.Target, etype, s.Note), nil
		case "http":
			return domain.MustNewHTTPEndpoint(s.Target, etype, s.UseProxy, spec.ProxyURL, s.Note), nil
		default:
			return domain.Endpoint{}, fmt.Errorf("unknown endpoint kind: %s", s.Kind)
		}
	}

	var vpn, direct, proxy []domain.Endpoint
	var err error

	for _, e := range spec.VPNEndpoints {
		ep, eerr := toEndpoint(e)
		if eerr != nil {
			err = fmt.Errorf("VPN endpoint %q: %w", e.Target, eerr)
			break
		}
		vpn = append(vpn, ep)
	}
	for _, e := range spec.DirectEndpoints {
		ep, eerr := toEndpoint(e)
		if eerr != nil {
			err = fmt.Errorf("Direct endpoint %q: %w", e.Target, eerr)
			break
		}
		direct = append(direct, ep)
	}
	for _, e := range spec.ProxyEndpoints {
		ep, eerr := toEndpoint(e)
		if eerr != nil {
			err = fmt.Errorf("Proxy endpoint %q: %w", e.Target, eerr)
			break
		}
		proxy = append(proxy, ep)
	}

	if err != nil {
		return domain.NetTestConfig{}, err
	}

	return domain.NewNetTestConfig(vpn, direct, proxy, spec.ProxyURL)
}
