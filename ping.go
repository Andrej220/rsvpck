package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func PingHostCmd(host string, attempts int, overallTimeout time.Duration) (bool, string, error) {
	if attempts < 1 {
		attempts = 1
	}
	if overallTimeout <= 0 {
		overallTimeout = 5 * time.Second
	}

	var args []string
	switch runtime.GOOS {
	case "windows":
		args = []string{"-n", fmt.Sprint(attempts), "-4", host}
	default: // linux, darwin, *bsd
		args = []string{"-c", fmt.Sprint(attempts), "-4", host}
	}

	ctx, cancel := context.WithTimeout(context.Background(), overallTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ping", args...)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return false, out.String(), fmt.Errorf("ping timed out after %s", overallTimeout)
	}

	if err != nil {
		text := out.String()
		if e := errb.String(); e != "" {
			text += "\n" + e
		}
		return false, text, nil
	}
	return true, out.String(), nil
}
