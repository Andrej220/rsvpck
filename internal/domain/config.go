package domain

import "errors"

type NetTestConfig struct {

	VPNEndpoints 		[]Endpoint
	DirectEndpoints 	[]Endpoint
	ProxyEndpoints 		[]Endpoint
	ProxyURL 			string
}

func NewNetTestConfig(
	vpnEndpoints 	[]Endpoint,
	directEndpoints []Endpoint,
	ProxyEndpoints 	[]Endpoint,
	proxyURL 		string,
) (NetTestConfig, error) {
	for _, ep := range vpnEndpoints {
		if ep.Type != EndpointTypeVPN  {
			return NetTestConfig{}, errors.New("all VPN endpoints must be of type VPN")
		}
		if ep.TargetType != TargetTypeTCP  && ep.TargetType != TargetTypeICMP{
			return NetTestConfig{}, errors.New("VPN endpoints must be TCP or ICMP")
		}
	}

	for _, ep := range directEndpoints {
		if ep.Type != EndpointTypePublic {
			return NetTestConfig{}, errors.New("direct endpoints must be of type Public")
		}
		switch ep.TargetType {
		case TargetTypeTCP, TargetTypeICMP, TargetTypeDNS, TargetTypeHTTP:
		default:
			return NetTestConfig{}, errors.New("direct endpoints must be TCP, ICMP, DNS, or HTTP")
		}	
	}
	for _, ep := range ProxyEndpoints{
		if ep.Type != EndpointTypePublic {
			return NetTestConfig{}, errors.New("proxy endpoint must be of type Public")
		}
		switch ep.TargetType {
		case TargetTypeICMP, TargetTypeDNS, TargetTypeHTTP:
		default:
			return NetTestConfig{}, errors.New("proxy endpoints must be ICMP or HTTP")
		}
	}

	return NetTestConfig{
		VPNEndpoints:       vpnEndpoints,
		DirectEndpoints: 	directEndpoints,
		ProxyEndpoints:  	ProxyEndpoints,
		ProxyURL:           proxyURL,
	}, nil
}


func (c NetTestConfig) HasVPNChecks() bool {
	return len(c.VPNEndpoints) > 0
}

func (c NetTestConfig) HasDirectChecks() bool {
	return len(c.DirectEndpoints) > 0
}

func (c NetTestConfig) HasProxyChecks() bool {
	return len(c.ProxyEndpoints) > 0
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
	if err != nil {
		panic("invalid HTTP endpoint: " + url + " - " + err.Error())
	}
	ep.RequiresProxy = overProxy
	ep.ProxyURL = proxyURL
	ep.Description = desc
	return ep
}

func MustNewICMPEndpoint(host string, typ EndpointType, description string) Endpoint {
	ep, err := NewICMPEndpoint(host, typ, description)
	if err != nil {
		panic("invalid ICMP endpoint: " + host + " - " + err.Error())
	}
	return ep
}