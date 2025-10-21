package domain

import (
	"time"
)

func AnalyzeConnectivity(probes []Probe, _ NetTestConfig) ConnectivityResult {
	r := ConnectivityResult{
		Probes:    probes,
		Timestamp: time.Now(),
	}
	r.DetermineMode()
	return r
}

