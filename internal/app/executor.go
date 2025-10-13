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
	policy      domain.ExecutionPolicy
}

func NewExecutor(
	tcpChecker 	domain.TCPChecker,
	dnsChecker 	domain.DNSChecker,
	httpChecker domain.HTTPChecker,
	icmpChecker domain.ICMPChecker,
	policy      domain.ExecutionPolicy,
) *Executor {
	return &Executor{
		tcpChecker:  tcpChecker,
		dnsChecker:  dnsChecker,
		httpChecker: httpChecker,
		icmpChecker: icmpChecker,
		policy:		 policy,
	}
}

func (e *Executor) Run(ctx context.Context, config domain.NetTestConfig) domain.ConnectivityResult {
    var probes []domain.Probe
	if e.policy == domain.PolicyOptimized{
		probes = append(probes, e.runOptimizedChecks(ctx, config.DirectEndpoints)...)
	} else{
		probes = append(probes, e.runEndpointCheck(ctx, config.DirectEndpoints)...)
	}
	probes = append(probes, e.runEndpointCheck(ctx,config.ProxyEndpoints)...)
	if e.policy == domain.PolicyOptimized{
		probes = append(probes, e.runOptimizedChecks(ctx, config.VPNEndpoints)...)
	} else {
		probes = append(probes, e.runEndpointCheck(ctx, config.VPNEndpoints)...)
	}

    return domain.AnalyzeConnectivity(probes, config)
}

func (e *Executor) runOptimizedChecks(ctx context.Context, endpoints []domain.Endpoint) []domain.Probe {
	var probes []domain.Probe
	var ipReachable bool

	var ipEndpoints, otherEndpoints []domain.Endpoint
	for _, ep := range endpoints {
		if ep.TargetType == domain.TargetTypeICMP  {
			ipEndpoints = append(ipEndpoints, ep)
		} else {
			otherEndpoints = append(otherEndpoints, ep)
		}
	}

	ipProbes := e.runEndpointCheck(ctx, ipEndpoints)
	probes = append(probes, ipProbes...)

	for _, p := range ipProbes {
		if p.IsSuccessful() {
			ipReachable = true
			break
		}
	}
	if ipReachable {
		otherProbes := e.runEndpointCheck(ctx, otherEndpoints)
		probes = append(probes, otherProbes...)
	}

	return probes
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