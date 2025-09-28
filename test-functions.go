package main

import (
	"time"
	"net/http"
	"net"
	"fmt"
	"net/url"
	"strings"
)

func testDNSResolution(domain string) NetTestResult {
    result := NetTestResult{TestName: "DNS Resolution", TestShortName: "DNS"}

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
        result.Details = fmt.Sprintf("TCP connect failed to %s", address)
        result.Error = err
		return result
	}
		
	_ = conn.Close()
	result.Status = StatusPass
	result.Details = fmt.Sprintf("TCP connect OK to %s", address)

    return result
}

func testInternetConnectivity(config *NetTestConfig) NetTestResult {
    result := NetTestResult{TestName: "Internet Connectivity", TestShortName: "Internet"}

    start := time.Now()
    conn, err := net.DialTimeout("tcp", internetConnectivityTestIP, config.Timeout)
	result.Latency = time.Since(start) 
    
    if err != nil {
        result.Status = StatusFail
        result.Details = "No internet connection"
        result.Error = err
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
		return result
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		result.Status = StatusPass
		result.Details = fmt.Sprintf("Connected successfully, %s", endpoint)
		return result
	} 

	result.Status = StatusWarning
	result.Details = fmt.Sprintf("%s responded with HTTP %d", endpoint,resp.StatusCode)
	return result
}

func testProxyHTTP(proxyAddr string, targetURL string, timeout time.Duration) NetTestResult {
    result := NetTestResult{
        TestName: fmt.Sprintf("Proxy via %s", proxyAddr),
		TestShortName: "Proxy",
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
    if err != nil {
        result.Status = StatusFail
        result.Details = fmt.Sprintf("Cannot reach %s via proxy %s", targetURL, proxyAddr)
        result.Error = err
        return result
    }
    defer resp.Body.Close()

    result.Latency = time.Since(start)

    if resp.StatusCode == http.StatusOK {
        result.Status = StatusPass
        result.Details = fmt.Sprintf("Proxy %s connected successfully to %s", proxyAddr, targetURL)
    } else {
        result.Status = StatusWarning
        result.Details = fmt.Sprintf("Proxy %s returned HTTP %d when connecting to %s", proxyAddr, resp.StatusCode, targetURL)
    }

    return result
}
