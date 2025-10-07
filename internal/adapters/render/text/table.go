// internal/adapters/render/text/table_renderer.go
package text

import (
	"fmt"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/azargarov/rsvpck/internal/domain"
	"io"
)

type TableRenderer struct{}

func NewTableRenderer() *TableRenderer {
	return &TableRenderer{}
}

func (tr *TableRenderer) Render(w io.Writer, result domain.ConnectivityResult) error {
	// Group and sort probes
	vpn, direct, proxy := tr.groupProbes(result.Probes)

	// Print overall status
	tr.printSummary(w,result)

	// Print sections
	if len(vpn) > 0 {
		fmt.Println("\nVPN Connectivity")
		tr.renderProbeTable(w, vpn)
	}

	if len(direct) > 0 {
		fmt.Println("\nDirect Internet")
		tr.renderProbeTable(w, direct)
	}

	if len(proxy) > 0 {
		fmt.Println("\nInternet via Proxy")
		tr.renderProbeTable(w, proxy)
	}

	return nil
}

func (tr *TableRenderer) printSummary(w io.Writer, result domain.ConnectivityResult) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	status := red("None")
	if result.IsConnected {
		status = green("Connected")
	}

	mode := tr.modeString(result.Mode)
	fmt.Fprintf(w,"%s • Mode: %s\n", status, mode)
	fmt.Fprintln(w,result.Summary)
}

func (tr *TableRenderer) modeString(mode domain.ConnectivityMode) string {
	switch mode {
	case domain.ModeDirect:
		return "Direct"
	case domain.ModeViaProxy:
		return "Via Proxy"
	case domain.ModeViaVPN:
		return "Via VPN"
	default:
		return "None"
	}
}

func (tr *TableRenderer) groupProbes(probes []domain.Probe) (vpn, direct, proxy []domain.Probe) {
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

func (tr *TableRenderer) renderProbeTable(w io.Writer, probes []domain.Probe) {
	table := tablewriter.NewTable(w,
        tablewriter.WithAlignment([]tw.Align{tw.AlignLeft, tw.AlignLeft, tw.AlignRight, tw.AlignLeft}),  
        tablewriter.WithRowAutoWrap(tw.WrapNormal),    
        tablewriter.WithHeaderAutoWrap(tw.WrapTruncate), 
        tablewriter.WithMaxWidth(400),                
    )

	table.Header([]string{"Test", "Status", "Latency", "Details"})

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, p := range probes {
		desc := p.Endpoint.Description
		if desc == "" {
			desc = fmt.Sprintf("%s (%s)", p.Endpoint.Target, p.Endpoint.TargetType.String())
		}

		statusStr := red("✗ Fail")
		if p.IsSuccessful() {
			statusStr = green("✓ Pass")
		}

		latencyStr := "-"
		if p.IsSuccessful() {
			latencyStr = fmt.Sprintf("%.2f ms", p.LatencyMs)
		}

		details := ""
		if !p.IsSuccessful() && p.Error != "" {
			details = tr.truncate(p.Error, 50)
		}

		table.Append([]string{desc, statusStr, latencyStr, details})
	}

	table.Render()
}

func (tr *TableRenderer) truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}