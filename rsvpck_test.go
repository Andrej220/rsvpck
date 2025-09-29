package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func mustStartTCPListener(t *testing.T) (net.Listener, string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen tcp: %v", err)
	}
	return ln, ln.Addr().String()
}

func startDummyTCPServer(t *testing.T) (addr string, closeFn func()) {
	ln, addr := mustStartTCPListener(t)
	done := make(chan struct{})
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-done:
					return
				default:
					return
				}
			}
			// simple echo/close
			_ = conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
			_, _ = conn.Write([]byte("ok"))
			_ = conn.Close()
		}
	}()
	return addr, func() { close(done); _ = ln.Close() }
}

func startHTTPProxy(t *testing.T) (proxyURL string, closeFn func()) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		target := r.URL.String()
		// Safety: only allow http target in tests
		if !strings.HasPrefix(target, "http://") {
			http.Error(w, "only http targets allowed in tests", http.StatusBadRequest)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), r.Method, target, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		req.Header = r.Header.Clone()
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		
		defer resp.Body.Close()
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	})
	srv := httptest.NewServer(handler)
	return srv.URL, srv.Close
}


func Test_testPortAvailability(t *testing.T) {
	okAddr, closeOK := startDummyTCPServer(t)
	defer closeOK()

	t.Run("ok", func(t *testing.T) {
		res := testPortAvailability(okAddr, 500*time.Millisecond)
		if res.Status != StatusPass {
			t.Fatalf("expected Pass, got %v details=%s err=%v", res.Status, res.Details, res.Error)
		}
		if res.Latency <= 0 {
			t.Fatalf("expected latency > 0")
		}
	})

	t.Run("bad address format -> Skipped", func(t *testing.T) {
		res := testPortAvailability("127.0.0.1", 200*time.Millisecond) // no port
		if res.Status != StatusSkipped {
			t.Fatalf("expected Skipped, got %v details=%s", res.Status, res.Details)
		}
	})

	t.Run("closed port -> Fail", func(t *testing.T) {
		res := testPortAvailability("127.0.0.1:1", 200*time.Millisecond)
		if res.Status != StatusFail {
			t.Fatalf("expected Fail, got %v details=%s", res.Status, res.Details)
		}
	})
}

func Test_testEndpoints(t *testing.T) {

	okSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer okSrv.Close()

	failSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer failSrv.Close()

	t.Run("OK -> Pass", func(t *testing.T) {
		res := testEndpoints(okSrv.URL, 1*time.Second)
		if res.Status != StatusPass {
			t.Fatalf("expected Pass, got %v details=%s err=%v", res.Status, res.Details, res.Error)
		}
	})

	t.Run("500 -> Warning", func(t *testing.T) {
		res := testEndpoints(failSrv.URL, 1*time.Second)
		if res.Status != StatusWarning {
			t.Fatalf("expected Warning, got %v details=%s", res.Status, res.Details)
		}
	})

	t.Run("unreachable -> Fail", func(t *testing.T) {
		res := testEndpoints("http://127.0.0.1:0", 300*time.Millisecond)
		if res.Status != StatusFail {
			t.Fatalf("expected Fail, got %v details=%s err=%v", res.Status, res.Details, res.Error)
		}
	})
}

func Test_testProxyHTTP(t *testing.T) {
	// target HTTP server (must be http:// for our tiny proxy)
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "proxied")
	}))
	defer target.Close()

	proxy, closeProxy := startHTTPProxy(t)
	defer closeProxy()

	proxyHost := strings.TrimPrefix(proxy, "http://")

	t.Run("proxy ok -> Pass", func(t *testing.T) {
		res := testProxyHTTP(proxyHost, target.URL, 2*time.Second)
		if res.Status != StatusPass {
			t.Fatalf("expected Pass, got %v details=%s err=%v", res.Status, res.Details, res.Error)
		}
	})
}

func Test_testDNSResolution(t *testing.T) {
	t.Run("localhost -> Pass", func(t *testing.T) {
		res := testDNSResolution("localhost")
		if res.Status != StatusPass {
			t.Fatalf("expected Pass, got %v details=%s", res.Status, res.Details)
		}
	})

	t.Run("invalid domain -> Fail", func(t *testing.T) {
		res := testDNSResolution("no.such.tld.invalid.")
		if res.Status != StatusFail {
			t.Fatalf("expected Fail, got %v details=%s", res.Status, res.Details)
		}
	})
}

func Test_testInternetConnectivity_OverrideTarget(t *testing.T) {
	addr, closeFn := startDummyTCPServer(t)
	defer closeFn()

	orig := internetConnectivityTestIP
	internetConnectivityTestIP = addr
	defer func() { internetConnectivityTestIP = orig }()

	cfg := &NetTestConfig{Timeout: 500 * time.Millisecond}
	res := testInternetConnectivity(cfg)
	if res.Status != StatusPass {
		t.Fatalf("expected Pass, got %v details=%s err=%v", res.Status, res.Details, res.Error)
	}

	internetConnectivityTestIP = "127.0.0.1:1"
	res = testInternetConnectivity(cfg)
	if res.Status != StatusFail {
		t.Fatalf("expected Fail, got %v details=%s err=%v", res.Status, res.Details, res.Error)
	}
}

func Test_getRoutePath_TimeoutDoesNotHang(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	done := make(chan NetTestResult, 1)
	go func() {
		done <- getRoutePath()
	}()

	select {
	case res := <-done:
		switch res.Status {
		case StatusUnknown, StatusFail, StatusPass, StatusWarning, StatusSkipped:
			// ok
		default:
			t.Fatalf("unexpected status: %v", res.Status)
		}
	case <-ctx.Done():
		t.Fatalf("getRoutePath timed out")
	}
}

// quick check that latency formatting never panics for 0/short durations
func Test_latencyToString_Smoke(t *testing.T) {
	got := latencyToString(0)
	if got == "" {
		t.Fatalf("latencyToString returned empty string for 0")
	}
	got = latencyToString(1234 * time.Microsecond)
	if got == "" {
		t.Fatalf("latencyToString returned empty string for short duration")
	}
}


func Test_Config_String_IncludesFields(t *testing.T) {
	cfg := &NetTestConfig{
		SiteID:   "site-123",
		HostName: "host-1",
		TestDate: time.Unix(0, 0).UTC(),
	}
	s := cfg.String()
	for _, want := range []string{"site-123", "host-1", "1970-01-01"} {
		if !strings.Contains(s, want) {
			t.Fatalf("config.String() missing %q in %q", want, s)
		}
	}
}

func Test_NetTestResult_String_Basic(t *testing.T) {
	r := &NetTestResult{
		TestName:     "X",
		TestShortName: "x",
		Status:       StatusPass,
		Details:      "ok",
		Latency:      10 * time.Millisecond,
	}
	s := r.String()
	for _, want := range []string{"X:", "ok", "latency"} {
		if !strings.Contains(s, want) {
			t.Fatalf("result.String() missing %q in %q", want, s)
		}
	}
}

func Test_PrintNetTestResult_Smoke(t *testing.T) {
	cfg := &NetTestConfig{SiteID: "s", HostName: "h", TestDate: time.Now()}
	data := []NetTestResult{
		{TestName: "A", TestShortName: "a", Status: StatusPass, Details: "ok"},
	}
	PrintNetTestResult(data, *cfg)
}
