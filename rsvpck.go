package main

import(
	"time"
    "fmt"
    "os"
)

const (
    timeout = 5 * time.Second
    applicationName = "RSvP connectivity checker"
)

var version = "dev"
var internetConnectivityTestIP = "8.8.8.8:53"

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

    for _, arg := range os.Args[1:] {
    	if arg == "-v" || arg == "--version" || arg == "-version" {
    		fmt.Printf("%s version %s\n", applicationName, version)
    		return
    	}
    }

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
