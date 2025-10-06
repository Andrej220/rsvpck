package icmp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/azargarov/rsvpck/internal/domain"
)

type Checker struct{}

var _ domain.ICMPChecker = (*Checker)(nil)

func (c *Checker)CheckPingWithContext(ctx context.Context, ep domain.Endpoint) domain.Probe{

	start := time.Now()
	ok, output, err := pingHostCmd(ctx, ep.Target, 1)
	latencyMs := time.Since(start).Seconds() * 1000

	if err != nil  || !ok {
		return domain.NewFailedProbe(
			ep,
			mapPingError(err, ctx.Err(), output),
			err,
		)
	}

	return domain.NewSuccessfulProbe(ep, latencyMs)
}

func pingHostCmd(ctx context.Context, host string, attempts int) (bool, string, error){
	if attempts < 1 {
		attempts = 1
	}

	var args []string
	switch runtime.GOOS {
	case "windows":
		args = []string{"-n", fmt.Sprint(attempts), "-4", host}
	default: // linux, darwin, *bsd
		args = []string{"-c", fmt.Sprint(attempts), "-4",host}
	}

	cmd := exec.CommandContext(ctx, "ping", args...)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	err := cmd.Run()

	output := out.String()
	if e := errb.String(); e != "" {
		output += "\n" + e
	}

	if ctx.Err() != nil {
		return false, output, ctx.Err()
	}

	// If the ping command failed to start (e.g. "ping: command not found")
	if err != nil {
		return false, output, err
	}

	success := isPingSuccessful(output, runtime.GOOS, attempts)

	return success, out.String(), nil
}

func isPingSuccessful(output, os string, attempts int) bool {
	outputLower := strings.ToLower(output)

	// On Linux/macOS: look for "0% packet loss" or received packets
	if os != "windows" {
		// Example: "5 packets transmitted, 5 received, 0% packet loss"
		if strings.Contains(outputLower, "received") &&
			!strings.Contains(outputLower, "100% packet loss") {
			return true
		}
		// Some systems: "5 packets transmitted, 0 received"
		if strings.Contains(outputLower, "0 received") {
			return false
		}
		return strings.Contains(outputLower, "bytes from")
	}

	// Windows: look for "Received = X" and "Lost = 0"
	if strings.Contains(outputLower, "received =") {
		return !strings.Contains(outputLower, "lost = "+fmt.Sprint(attempts))
	}
	// Or: "Reply from ..."
	return strings.Contains(outputLower, "reply from")
}

func mapPingError(err error, contextErr error, output string) domain.Status {
	// 1. Context cancellation/timeout
	if contextErr != nil {
		if errors.Is(contextErr, context.DeadlineExceeded) ||
			errors.Is(contextErr, context.Canceled) {
			return domain.StatusTimeout
		}
	}

	// 2. Parse output for DNS-like errors
	outputLower := strings.ToLower(output)
	if containsAny(outputLower,
		"unknown host",
		"name or service not known",
		"nodename nor servname provided",
		"cannot resolve",
		"no address associated",
	) {
		return domain.StatusDNSFailure
	}

	// 3. Network unreachable / timeout
	if containsAny(outputLower,
		"network is unreachable",
		"host is unreachable",
		"operation timed out",
		"request timeout",
		"100% packet loss",
		"destination host unreachable",
	) {
		return domain.StatusTimeout
	}

	// 4. General failure
	return domain.StatusInvalid
}

func containsAny(s string, substrs ...string) bool {
	sLower := strings.ToLower(s)
	for _, substr := range substrs {
		if substr != "" && strings.Contains(sLower, strings.ToLower(substr)) {
			return true
		}
	}
	return false
}