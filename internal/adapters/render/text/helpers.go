package text

import (
	"fmt"
	"github.com/azargarov/rsvpck/internal/domain"
	"github.com/fatih/color"
	"io"
)

func printSummary(w io.Writer, result domain.ConnectivityResult) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	status := red("None")
	if result.IsConnected {
		status = green("Connected")
	}

	mode := modeString(result.Mode)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%s • Mode: %s\n", status, mode)
	fmt.Fprintln(w, "")
}

func modeString(mode domain.ConnectivityMode) string {
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
