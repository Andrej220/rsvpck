package main

import(
	"time"
    "fmt"
    "os"
)

const (
    timeout = 5 * time.Second
    applicationName = "RSvP connectivity checker"
    timeoutMarker = time.Hour
)

var version = "dev"
var internetConnectivityTestIP = "8.8.8.8:53"

var endpoints = []string{
                "https://insite-eu.gehealthcare.com",
                //"https://insite.gehealthcare.com",
}

var proxies = []string{
                "54.154.45.26:443",
                "82.136.152.78:8002",
                "10.25.0.20:8080",
                "152.2.1.251:8002",
}

var toping = []string{  "150.2.101.89",
                        "82.136.152.65",
                        "8.8.8.8",
                    }

func RunRSVPDiagnostics(config *NetTestConfig) []NetTestResult {
    var results []NetTestResult

    fmt.Println("Checking Internet connectivity")

    results = append(results, testInternetConnectivity(config))
    for _, s := range(config.CheckEndpoints){
        results = append(results, testDNSResolution(s))
    }

    fmt.Println("Checking endpoints")

    for _, s := range(config.CheckEndpoints){
        results = append(results, testEndpoints(s, config.Timeout))
    }
    for _, px := range(config.CheckProxies){
        results = append(results, testPortAvailability(px,config.Timeout))
    }

    fmt.Println("Testing proxies")


    for _, px := range(config.CheckProxies){
        for _, s := range(config.CheckEndpoints){
            results = append(results, testProxyHTTP(px, s, config.Timeout))
        }
    }

    return results
}

func collectHostData()[]string{

    var results []string
    results = append(results, getRoutePath().Details)
    return results
}

func pingProxies( hosts []string) []NetTestResult{
    fmt.Println("Pinging...")
    var results []NetTestResult
    for _, h := range(hosts){
        results = append(results, pingAProxy(h))
    }
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
    PrintNetTestResult(results, *config)
    results = pingProxies(toping)
    PrintNetTestResult(results, *config)

    hostData := collectHostData()
    for _,s := range(hostData){
        printText(s,"Routing table")
    }

    
    tls, err := CheckTLS("insite-eu.gehealthcare.com:443","insite-eu.gehealthcare.com", 5* time.Second)
    if err != nil {
        fmt.Println("Failed to read TLS certificates")
    } else{
        printText(tls,"Certificates")
    }
    
    
}
