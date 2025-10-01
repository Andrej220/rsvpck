package main

import (
	"fmt"
	"time"
	"net/url"
	"strings"

)

func latencyToString(d time.Duration) string {
    switch {
    case d == 0:
        return "-" 
    case d >= timeoutMarker:
        return "timeout"
    default:
        ms := float64(d.Microseconds()) / 1000.0
        if ms < 1000 {
            return fmt.Sprintf("%.2f ms", ms)
        }
        return fmt.Sprintf("%.0f ms", ms)
    }
}

func indentMultiline(s string) string {
    s = strings.TrimSpace(s)
    if s == "" {
        return "-"
    }
    return strings.ReplaceAll(s, "\n", "\n\t")
}

func hostFromEndpoint(ep string) (string, error){
	ep = strings.TrimSpace(ep)
	if ep == "" {
		return "", fmt.Errorf("empty endpoint")
	}

    // url.Parse needs a scheme, fake one if missing
    if !strings.Contains(ep, "://") {
        ep = "http://" + ep
    }

    u, err := url.Parse(ep)
    if err != nil {
        return ep, err 
    }
    return u.Hostname(), nil
}

func btoi(b bool) int { if b { return 1 }; return 0 }