package text

import (
	"fmt"
	"io"
	"sort"
	"github.com/azargarov/rsvpck/internal/domain"
)

type TextRenderer struct{
	conf *RenderConfig
}

func NewRenderer(conf *RenderConfig) *TextRenderer {
	return &TextRenderer{conf: conf}
}

func (r *TextRenderer) Render(w io.Writer, result domain.ConnectivityResult) error {

	vpnProbes, directProbes, proxyProbes := r.groupProbes(result.Probes)
	
	if len(vpnProbes) > 0 {
		fmt.Fprintln(w, "VPN Connectivity:")
		r.renderProbeList(w, vpnProbes)
		fmt.Fprintln(w)
	}
	
	if len(directProbes) > 0 {
		fmt.Fprintln(w, "Direct Internet:")
		r.renderProbeList(w, directProbes)
		fmt.Fprintln(w)
	}
	
	if len(proxyProbes) > 0 {
		fmt.Fprintln(w, "Internet via Proxy:")
		r.renderProbeList(w, proxyProbes)
		fmt.Fprintln(w)
	}
	
	printSummary(w, result, r.conf)

	return nil
}

func (r *TextRenderer) groupProbes(probes []domain.Probe) (vpn, direct, proxy []domain.Probe) {
	for _, p := range probes {
		if p.Endpoint.Type == domain.EndpointTypeVPN {
			vpn = append(vpn, p)
		} else if p.Endpoint.Proxy.MustUseProxy(){
			proxy = append(proxy, p)
		} else {
			direct = append(direct, p)
		}
	}
	return
}

func (r *TextRenderer) renderProbeList(w io.Writer, probes []domain.Probe) {
	sort.Slice(probes, func(i, j int) bool {
		if probes[i].IsSuccessful() != probes[j].IsSuccessful() {
			return probes[i].IsSuccessful()
		}
		return probes[i].Endpoint.Target < probes[j].Endpoint.Target
	})

	for _, p := range probes {
		statusIcon := r.conf.FailSym
		if p.IsSuccessful() {
			statusIcon = r.conf.OkSym
		}

		desc := p.Endpoint.Description
		if desc == "" {
			desc = p.Endpoint.Target
		}

		if p.IsSuccessful() {
			fmt.Fprintf(w, "\t%s %-40s [%.2f ms]\n", statusIcon, desc, p.LatencyMs)
		} else {
			errorMsg := truncateError(p.Error, maxCharPerError)
			fmt.Fprintf(w, "\t%s %-40s %s\n", statusIcon, desc, errorMsg)
		}
	}
}
