package domain

import "errors"

type NetTestConfig struct {

	VPNEndpoints 		[]Endpoint
	DirectEndpoints 	[]Endpoint
	HTTPEndpoints 		[]Endpoint
	DNSEndpoints 		[]Endpoint
	ProxyURL 			string
}

func NewNetTestConfig(
	vpnEndpoints 	[]Endpoint,
	directEndpoints []Endpoint,
	HTTPEndpoints 	[]Endpoint,
	DnsEndpoints    []Endpoint,
	proxyURL 		string,
) (NetTestConfig, error) {
	for _, ep := range vpnEndpoints {
		if ep.Type != EndpointTypeVPN  {
			return NetTestConfig{}, errors.New("all VPN endpoints must be of type VPN")
		}
		if ep.TargetType != TargetTypeTCP  && ep.TargetType != TargetTypeICMP{
			return NetTestConfig{}, errors.New("VPN endpoints must be TCP (host:port) or ICMP")
		}
	}

	for _, ep := range directEndpoints {
		if ep.Type != EndpointTypePublic {
			return NetTestConfig{}, errors.New("direct endpoints must be of type Public")
		}
		if (ep.TargetType != TargetTypeTCP && ep.TargetType != TargetTypeICMP) {
			return NetTestConfig{}, errors.New("direct endpoints must be TCP (host:port) or ICMP type")
		}
	}
	for _, ep := range HTTPEndpoints{
		if ep.Type != EndpointTypePublic {
			return NetTestConfig{}, errors.New("proxy endpoint must be of type Public")
		}
		if ep.TargetType != TargetTypeHTTP {
			return NetTestConfig{}, errors.New("proxy endpoint must be HTTP (URL)")
		}
	}

	return NetTestConfig{
		VPNEndpoints:       vpnEndpoints,
		DirectEndpoints: 	directEndpoints,
		HTTPEndpoints:  	HTTPEndpoints,
		DNSEndpoints:       DnsEndpoints,
		ProxyURL:           proxyURL,
	}, nil
}


func (c NetTestConfig) HasVPNChecks() bool {
	return len(c.VPNEndpoints) > 0
}

func (c NetTestConfig) HasDirectChecks() bool {
	return len(c.DirectEndpoints) > 0
}

func (c NetTestConfig) HasProxy() bool {
	return len(c.HTTPEndpoints) > 0
}

func (c NetTestConfig) HasDNSCheck() bool {
	return len(c.DNSEndpoints) > 0
}

func (c NetTestConfig) HasICMPChecks() bool {
	for _, ep := range c.DirectEndpoints {
		if ep.TargetType == TargetTypeICMP {
			return true
		}
	}
	for _, ep := range c.VPNEndpoints {
		if ep.TargetType == TargetTypeICMP {
			return true
		}
	}
	return false
}

func MustNewTCPEndpoint(hostPort string, typ EndpointType, desc string) Endpoint {
	ep, err := NewTCPEndpoint(hostPort, typ, desc)
	if err != nil {
		panic("invalid endpoint: " + hostPort + " - " + err.Error())
	}
	return ep
}

func MustNewDNSEndpoint(host string, typ EndpointType, desc string) Endpoint {
	ep, err := NewDNSEndpoint(host, typ, desc)
	if err != nil {
		panic("invalid DNS endpoint: " + host + " - " + err.Error())
	}
	return ep
}

func MustNewHTTPEndpoint(url string, typ EndpointType, overProxy bool, proxyURL string ,desc string) Endpoint {
	ep, err := NewHTTPEndpoint(url, typ, desc)
	ep.RequiresProxy = overProxy
	ep.ProxyURL = proxyURL
	ep.Description = desc
	if err != nil {
		panic("invalid HTTP endpoint: " + url + " - " + err.Error())
	}
	return ep
}

func MustNewICMPEndpoint(host string, typ EndpointType, description string) Endpoint {
	ep, err := NewICMPEndpoint(host, typ, description)
	if err != nil {
		panic("invalid ICMP endpoint: " + host + " - " + err.Error())
	}
	return ep
}