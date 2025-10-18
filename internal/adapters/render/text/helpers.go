package text

import (
	"fmt"
	"github.com/azargarov/rsvpck/internal/domain"
	"io"
	"strings"
)

func printSummary(w io.Writer, result domain.ConnectivityResult, conf *RenderConfig) {

	status := conf.Red("None")
	if result.IsConnected {
		status = conf.Green("Connected")
	}

	mode := modeString(result.Mode)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%s > Mode: %s\n", status, mode)
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

func truncateError(msg string, maxLen int) string {
	if strings.Contains(msg, "UTF-8") || strings.Contains(msg, "UTF8") {
		return truncateErrorUTF8(msg, maxLen)
	} 
	return truncateErrorASCII(msg, maxLen)
}

func truncateErrorASCII(msg string, maxLen int) string {
	if msg == "" {
		return ""
	}
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen-3] + "..."
}

func truncateErrorUTF8(msg string, max int) string {
	if len(msg) <= max {
		return msg
	}
	rn := []rune(msg)
	if len(rn) <= max {
		return msg
	}
	if max < 3 {
		return string(rn[:max])
	}
	return string(rn[:max-3]) + "..."
}
