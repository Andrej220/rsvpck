package text

import (
	"fmt"
	"io"
	"sort"
	"github.com/fatih/color"
	"github.com/azargarov/rsvpck/internal/domain"
)
type TextRenderer struct{}

func NewRenderer() *TextRenderer {
	return &TextRenderer{}
}

func (r *TextRenderer) Render(w io.Writer, result domain.ConnectivityResult) error {
	
	printSummary(w,result)

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

	return nil
}

func (r *TextRenderer) groupProbes(probes []domain.Probe) (vpn, direct, proxy []domain.Probe) {
	for _, p := range probes {
		if p.Endpoint.Type == domain.EndpointTypeVPN {
			vpn = append(vpn, p)
		} else if p.Endpoint.RequiresProxy {
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

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, p := range probes {
		statusIcon := green("✓ ")
		if !p.IsSuccessful() {
			statusIcon = red("✗ ")
		}

		desc := p.Endpoint.Description
		if desc == "" {
			desc = p.Endpoint.Target
		}

		if p.IsSuccessful() {
			fmt.Fprintf(w, "\t%s %-40s [%.2f ms]\n", statusIcon, desc, p.LatencyMs)
		} else {
			errorMsg := r.truncateError(p.Error, 50)
			fmt.Fprintf(w, "\t%s %-40s %s\n", statusIcon, desc, errorMsg)
		}
	}
}

func (r *TextRenderer) truncateError(msg string, maxLen int) string {
	if msg == "" {
		return ""
	}
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}