package main

import (
	"fmt"
	"time"
	"net/url"
	"strings"

)

func latencyToString(l time.Duration) string{
	latencyMs := float64(l.Microseconds()) / 1000.0
	return fmt.Sprintf("%.2f ms", latencyMs)
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