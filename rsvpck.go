package main

import(
	"time"
)

const (
	internetConnectivityTestIP = "8.8.8.8:53"
    timeout = 5 * time.Second
)

var endpoints = []string{
                "https://insite-eu.gehealthcare.com",
                "https://insite.gehealthcare.com",
}

var proxies = []string{
                ":443",
                ":8002",
}

func RunRSVPDiagnostics(config *NetTestConfig) []NetTestResult {
    var results []NetTestResult
    
    results = append(results, testInternetConnectivity(config))
    for _, s := range(config.CheckEndpoints){
        results = append(results, testDNSResolution(s))
        results = append(results, testEndpoints(s, config.Timeout))
    }
    for _, px := range(config.CheckProxies){
        results = append(results, testPortAvailability(px,config.Timeout))
        for _, s := range(config.CheckEndpoints){
            results = append(results, testProxyHTTP(px, s, config.Timeout))
        }
    }
    results = append(results, getRoutePath())

    return results
}

func main(){

    config := &NetTestConfig{
            SiteID:  getMachineUUID(),
            HostName: getHostName(),
            TestDate: time.Now(),
            Timeout: timeout,
            CheckProxies: proxies,
            CheckEndpoints: endpoints, 
    }  
    results := RunRSVPDiagnostics(config)
    PrintNetTestResult(results, config)
}
