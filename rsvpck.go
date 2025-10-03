package main

import(
	"time"
    "fmt"
    "os"
)

const (
    timeout = 2 * time.Second
    applicationName = "RSvP connectivity checker"
    timeoutMarker = time.Hour
)

var version = "dev"
var internetConnectivityTestIP = "google.com:443"//"8.8.8.8:53"
var internetIP = []string {"8.8.8.8","1.1.1.1", "google.com"}

var endpoints = []string{
                "https://insite-eu.gehealthcare.com",
                //"https://insite.gehealthcare.com",
}
var internetProxy = []string{
                            "54.154.45.26:443",
                            }

var proxies = []string{
                "82.136.152.78:8002",
                "10.25.0.20:8080",
                "150.2.1.251:8002",
}

var proxyToPing = []string{  "150.2.101.89",
                        "82.136.152.65",
                    }

func RunInternetTest(config *NetTestConfig) []NetTestResult {
    var results []NetTestResult

    //fmt.Println("Internet connectivity")

    for _, h := range(internetIP){
        results = append(results, pingAProxy(h))
    }

    results = append(results, testInternetConnectivity(config))
    for _, s := range(config.CheckEndpoints){
        results = append(results, testDNSResolution(s))
    }

    for _, s := range(config.CheckEndpoints){
        results = append(results, testEndpoints(s, config.Timeout))
    }
    for _, px := range(internetProxy){
        results = append(results, testPortAvailability(px,config.Timeout))
    }

    for _, px := range(internetProxy){
        for _, s := range(config.CheckEndpoints){
            results = append(results, testProxyHTTP(px, s, config.Timeout))
        }
    }
 

    return results
}

func runVPNtest(hosts []string)[]NetTestResult{
    //fmt.Println("VPN")
    var results []NetTestResult
    for _, h := range(hosts){
        results = append(results, pingAProxy(h))
    }
    return results
}

func runProxyTest(config *NetTestConfig)[]NetTestResult{
    var results []NetTestResult
    //mt.Println("Proxy")

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
    fmt.Printf("%s, Version: %s\n", applicationName, version)
    fmt.Printf("SID: %s\n\n", config.SiteID)

    results := RunInternetTest(config)
    PrintNetTestResult("Internet",results, *config)

    results = runProxyTest(config)
    PrintNetTestResult("Proxy",results, *config) 

    results = runVPNtest(proxyToPing)
    PrintNetTestResult("VPN",results, *config)

    hostData := collectHostData()
    for _,s := range(hostData){
        printText(s,"Default route")
    }

    
    tls, err := CheckTLS("insite-eu.gehealthcare.com:443","insite-eu.gehealthcare.com", 5* time.Second)
    if err != nil {
        fmt.Println("Failed to read TLS certificates")
    } else{
        printText(tls,"Certificate")
    }
    
    
}
