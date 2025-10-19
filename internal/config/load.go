package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/azargarov/rsvpck/internal/domain"
	"gopkg.in/yaml.v3"
)

type FileSpec struct {
	ProxyURL        string         `json:"proxyURL"        yaml:"proxyURL"`
	VPNIPs			[]string	   `json:"vpnIPs"          yaml:"vpnIPs"`
	VPNEndpoints    []EndpointSpec `json:"vpnEndpoints"    yaml:"vpnEndpoints"`
	DirectEndpoints []EndpointSpec `json:"directEndpoints" yaml:"directEndpoints"`
	ProxyEndpoints  []EndpointSpec `json:"proxyEndpoints"  yaml:"proxyEndpoints"`
}

type EndpointSpec struct {
	Target   string `json:"target"   yaml:"target"`   
	Type     string `json:"type"     yaml:"type"`     
	Kind     string `json:"kind"     yaml:"kind"`     
	Note     string `json:"note"     yaml:"note"`     
	UseProxy bool   `json:"useProxy" yaml:"useProxy"` 
}

func LoadFromFile(path string) (domain.NetTestConfig, error) {
	if path == "" {
		return domain.NetTestConfig{}, errors.New("empty path")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return domain.NetTestConfig{}, err
	}
	return parseConfigBytes(raw, filepath.Ext(path))
}

func parseConfigBytes(b []byte, ext string) (domain.NetTestConfig, error) {
	ext = strings.ToLower(ext)
	var spec FileSpec
	var err error

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(b, &spec)
		if err != nil {
			if jerr := json.Unmarshal(b, &spec); jerr == nil {
				err = nil
			}
		}
	default:
		err = json.Unmarshal(b, &spec)
		if err != nil {
			if yerr := yaml.Unmarshal(b, &spec); yerr == nil {
				err = nil
			}
		}
	}
	if err != nil {
		return domain.NetTestConfig{}, fmt.Errorf("invalid config: %w", err)
	}

	return specToDomain(spec)
}

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
			err = fmt.Errorf("direct endpoint %q: %w", e.Target, eerr)
			break
		}
		direct = append(direct, ep)
	}
	for _, e := range spec.ProxyEndpoints {
		ep, eerr := toEndpoint(e)
		if eerr != nil {
			err = fmt.Errorf("proxy endpoint %q: %w", e.Target, eerr)
			break
		}
		proxy = append(proxy, ep)
	}
	
	if err != nil {
		return domain.NetTestConfig{}, err
	}

	return domain.NewNetTestConfig(vpn, direct, proxy, spec.ProxyURL, spec.VPNIPs)
}
