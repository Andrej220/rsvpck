package text

import (
	"fmt"
	"io"

	"github.com/azargarov/rsvpck/internal/domain"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

const maxTableWidth = 400

type TableRenderer struct{
	conf *RenderConfig
}

func NewTableRenderer(conf *RenderConfig) *TableRenderer {
	return &TableRenderer{conf: conf}
}

func (tr *TableRenderer) Render(w io.Writer, result domain.ConnectivityResult) error {
	// Group and sort probes
	vpn, direct, proxy := tr.groupProbes(result.Probes)

	// Print overall status
	
	// Print sections
	if len(vpn) > 0 {
		tr.renderProbeTable(w, vpn, "VPN Connectivity")
	}
	
	if len(direct) > 0 {
		tr.renderProbeTable(w, direct, "Direct Internet")
	}
	
	if len(proxy) > 0 {
		tr.renderProbeTable(w, proxy, "Internet via Proxy")
	}
	
	printSummary(w, result, tr.conf)
	
	return nil
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

func (tr *TableRenderer) renderProbeTable(w io.Writer, probes []domain.Probe, name string) {
	table := tablewriter.NewTable(w,
		tablewriter.WithAlignment([]tw.Align{tw.AlignLeft, tw.AlignLeft, tw.AlignRight, tw.AlignLeft}),
		tablewriter.WithRowAutoWrap(tw.WrapNormal),
		tablewriter.WithHeaderAutoWrap(tw.WrapTruncate),
		tablewriter.WithMaxWidth(maxTableWidth),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: tr.conf.TableSymbols})),
	)

	table.Header([]string{name, "Status", "Latency", "Details"})

	for _, p := range probes {
		desc := p.Endpoint.Description
		if desc == "" {
			desc = fmt.Sprintf("%s (%s)", p.Endpoint.Target, p.Endpoint.TargetType.String())
		}

		statusStr := tr.conf.FailSym + " Fail"
		if p.IsSuccessful() {
			statusStr = tr.conf.OkSym + " Pass"
		}

		latencyStr := "-"
		if p.IsSuccessful() {
			latencyStr = fmt.Sprintf("%.2f ms", p.LatencyMs)
		}

		details := ""
		if !p.IsSuccessful() && p.Error != "" {
			details = truncateError(p.Error, maxCharPerError)
		}

		table.Append([]string{desc, statusStr, latencyStr, details})
	}

	table.Render()
}

func getTableBorders() *tw.SymbolCustom {

	nature := tw.NewSymbolCustom("Nature").
		WithRow("─").
		WithColumn("│").
		WithTopLeft("┌").WithTopMid("┬").WithTopRight("┐").
		WithMidLeft("├").WithCenter("┼").WithMidRight("┤").
		WithBottomLeft("└").WithBottomMid("┴").WithBottomRight("┘")

	ascii := tw.NewSymbolCustom("ASCII").
		WithRow("-").
		WithColumn("|").
		WithTopLeft("+").WithTopMid("+").WithTopRight("+").
		WithMidLeft("+").WithCenter("+").WithMidRight("+").
		WithBottomLeft("+").WithBottomMid("+").WithBottomRight("+")

	sym := ascii
	if unicodeSupported {
		sym = nature
	}
	return sym
}
