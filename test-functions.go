package main

import (
	"time"
	"net/http"
	"net"
	"fmt"
	"net/url"
	"strings"
	"crypto/tls"
)

func testDNSResolution(domain string) NetTestResult {
    result := NetTestResult{TestName: "DNS Resolution", TestShortName: "DNS Resolution"}

	domain, err := hostFromEndpoint(domain)
	if err != nil{
		result.Status = StatusUnknown
		result.Details = fmt.Sprintf("Could not parse domain, %s", domain)
		return result
	}

	start := time.Now()
    addrs, err := net.LookupHost(domain)
	result.Latency = time.Since(start)

    if err != nil {
        result.Status = StatusFail
        result.Details = fmt.Sprintf("%s: UNRESOLVABLE(%v)", domain, err)
    } else {
		result.Status = StatusPass
        result.Details = fmt.Sprintf("%s: %v", domain, addrs)
    }
    
    return result
}

func pingAProxy(host string) NetTestResult {
    result := NetTestResult{TestName: "Ping", TestShortName: "Ping"}

	start := time.Now()
    ok, _, err := PingHostCmd(host, 1, 5 * time.Second)
	result.Latency = time.Since(start)
	if err != nil{
		result.Status = StatusUnknown
        result.Details = fmt.Sprintf("%s ", host)
		result.Latency = 0
	}

    if !ok {
        result.Status = StatusFail
        result.Details = fmt.Sprintf("%s ", host)
		result.Latency = timeoutMarker
    } else {
		result.Status = StatusPass
        result.Details = fmt.Sprintf("%s ", host)
		//result.Latency = 0
    }
    
    return result
}


func checkCertificates(host string)string{
		client := &http.Client{
		Timeout: 5 * time.Second, // total request timeout
	}

	resp, err := client.Get(host)
	if err != nil {
		return fmt.Sprint("TLS/HTTP check failed:", err)
	}
	defer resp.Body.Close()

	return fmt.Sprint("OK, certificate is valid. HTTP status:", resp.Status)
}

func CheckTLS(addr, serverName string, timeout time.Duration) (string, error) {
	dialer := &net.Dialer{Timeout: timeout}

	tlsCfg := &tls.Config{
		ServerName: serverName, // important for hostname verification & SNI
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsCfg)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var sb strings.Builder
	state := conn.ConnectionState()

	for i, cert := range state.PeerCertificates {
		sb.WriteString(fmt.Sprintf("Certificate #%d:\n", i))
		sb.WriteString(fmt.Sprintf("   Subject    : %s\n", cert.Subject))
		sb.WriteString(fmt.Sprintf("   Issuer     : %s\n", cert.Issuer))
		sb.WriteString(fmt.Sprintf("   Valid from : %s\n", cert.NotBefore.Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("   Valid until: %s\n\n", cert.NotAfter.Format("2006-01-02")))
	}

	return sb.String(), nil
}

// Address must include port
func testPortAvailability(address string, timeout time.Duration) NetTestResult {
    result := NetTestResult{TestName: "Ports availability", TestShortName: "Proxy precheck"}
	if !strings.Contains(address, ":") { 
		result.Status = StatusSkipped
		result.Details = fmt.Sprintf("Failed address string: %s", address)
		return result
	}

	d := net.Dialer{Timeout: timeout, KeepAlive: 15 * time.Second}

	start := time.Now()
	conn, err := d.Dial("tcp", address)
	result.Latency = time.Since(start)

	if err != nil{
        result.Status = StatusFail
        result.Details = fmt.Sprintf("%s", address)
        result.Error = err
		result.Latency = timeoutMarker
		return result
	}
		
	_ = conn.Close()
	result.Status = StatusPass
	result.Details = fmt.Sprintf("%s", address)

    return result
}

func testInternetConnectivity(config *NetTestConfig) NetTestResult {
    result := NetTestResult{TestName: "HTTPS", 
	                        TestShortName: "HTTPS: " + internetConnectivityTestIP}

    start := time.Now()
    conn, err := net.DialTimeout("tcp", internetConnectivityTestIP, config.Timeout)
	result.Latency = time.Since(start) 
    
    if err != nil {
        result.Status = StatusFail
        result.Details = "No internet connection"
        result.Error = err
		result.Latency = timeoutMarker
		return result
    } 

    _ = conn.Close()
    result.Status = StatusPass
    result.Details = "Internet connection available"
    
    return result
}

func testEndpoints( endpoint string, timeout time.Duration) NetTestResult {
result := NetTestResult{
		TestName: "GEHC Connectivity",
		TestShortName: "Endpoint",
	}

	client := &http.Client{Timeout: timeout}
	start := time.Now()
	req, err := http.NewRequest("GET", endpoint, nil)
	
	if err != nil {
		result.Status = StatusFail
		result.Details = "Failed to create request"
		result.Error = err
		return result
	}

	resp, err := client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Status = StatusFail
		result.Details = fmt.Sprintf("Cannot reach %s", endpoint)
		result.Error = err
		result.Latency = timeoutMarker
		return result
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		result.Status = StatusPass
		result.Details = fmt.Sprintf("%s", endpoint)
		return result
	} 

	result.Status = StatusWarning
	result.Details = fmt.Sprintf("%s responded with HTTP %d", endpoint,resp.StatusCode)
	return result
}

func testProxyHTTP(proxyAddr string, targetURL string, timeout time.Duration) NetTestResult {
    result := NetTestResult{
        TestName: fmt.Sprintf("Proxy via %s", proxyAddr),
		TestShortName: "Proxy: " + proxyAddr,
    }

    proxyURL, err := url.Parse("http://" + proxyAddr)
    if err != nil {
        result.Status = StatusFail
        result.Details = "Invalid proxy address"
        result.Error = err
        return result
    }

    transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
    client := &http.Client{Transport: transport, Timeout: timeout}

    start := time.Now()
    resp, err := client.Get(targetURL) 
	result.Latency = time.Since(start)

    if err != nil {
        result.Status = StatusFail
        result.Details = fmt.Sprintf("%s", targetURL)
        result.Error = err
		result.Latency = timeoutMarker
        return result
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        result.Status = StatusPass
        result.Details = fmt.Sprintf("%s", targetURL)
    } else {
        result.Status = StatusWarning
        result.Details = fmt.Sprintf("%s returned HTTP %d when connecting to %s", proxyAddr, resp.StatusCode, targetURL)
    }

    return result
}
