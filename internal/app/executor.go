package app

import (
	"github.com/azargarov/rsvpck/internal/domain"
	"context"
)

type Executor struct {
	tcpChecker  domain.TCPChecker
	dnsChecker  domain.DNSChecker
	httpChecker domain.HTTPChecker
	icmpChecker domain.ICMPChecker
}

func NewExecutor(
	tcpChecker 	domain.TCPChecker,
	dnsChecker 	domain.DNSChecker,
	httpChecker domain.HTTPChecker,
	icmpChecker domain.ICMPChecker,
) *Executor {
	return &Executor{
		tcpChecker:  tcpChecker,
		dnsChecker:  dnsChecker,
		httpChecker: httpChecker,
		icmpChecker: icmpChecker,
	}
}

func (e *Executor) Run(ctx context.Context, config domain.NetTestConfig) domain.ConnectivityResult {
    var probes []domain.Probe

	probes = append(probes, e.runEndpointCheck(ctx, config.VPNEndpoints)...)
	probes = append(probes, e.runEndpointCheck(ctx, config.DirectEndpoints)...)
	probes = append(probes, e.runEndpointCheck(ctx,config.ProxyEndpoints)...)

    return domain.AnalyzeConnectivity(probes, config)
}

func (e Executor)runEndpointCheck(ctx context.Context, endpoints []domain.Endpoint) []domain.Probe{
	var probes []domain.Probe
	for _, ep := range endpoints {
		var probe domain.Probe

		switch ep.GetTargetType(){
		case domain.TargetTypeICMP:
			probe =e.icmpChecker.CheckPingWithContext(ctx, ep)
		case domain.TargetTypeTCP:
			probe = e.tcpChecker.CheckWithContext(ctx, ep)
		case domain.TargetTypeHTTP:
			if ep.RequiresProxy{
				probe = e.httpChecker.CheckViaProxyWithContext(ctx, ep, ep.ProxyURL)
			} else{
				probe = e.httpChecker.CheckWithContext(ctx, ep)
			}
		case domain.TargetTypeDNS:
			probe = e.dnsChecker.CheckWithContext(ctx, ep)
		}

		probes = append(probes, probe)
	}
	return probes
}