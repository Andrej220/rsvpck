package main

import (
	"github.com/azargarov/rsvpck/internal/domain"
	"github.com/azargarov/rsvpck/internal/app"
	"github.com/azargarov/rsvpck/internal/adapters/tcp"
	"github.com/azargarov/rsvpck/internal/adapters/dns"
	"github.com/azargarov/rsvpck/internal/adapters/http"
	"github.com/azargarov/rsvpck/internal/adapters/icmp"
	"github.com/azargarov/rsvpck/internal/adapters/render/text"
	"github.com/azargarov/rsvpck/internal/adapters/hostinfo"
	"github.com/azargarov/rsvpck/internal/adapters/httpx"
	"github.com/azargarov/go-utils/autostr"

	"fmt"
	"time"
	"context"
	"os"
)

func main(){

	printHeader()
	rsvpConf := parseFlagsToConfig()

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	h := hostinfo.GetCRMInfo(ctx)
	autostrCfg := autostr.Config{Separator: autostr.Ptr("\n"), FieldValueSeparator: autostr.Ptr(" : "), PrettyPrint: true}

	text.PrintBlock(os.Stdout,"SYSTEM INFORMATION", autostr.String(h, autostrCfg))
	h.TLSCert, err = httpx.GetCertificates(ctx,"insite-eu.gehealthcare.com:443","insite-eu.gehealthcare.com")
	
	if err == nil{
		text.PrintList(os.Stdout, "TLS certificates, eu-insite.gehealthcare.com\n",h.TLSCert)
	}else{
		fmt.Println("Failed to fetch certificates")
	}

	tcpChecker := &tcp.Checker{}
	dnsChecker := &dns.Checker{}
	httpChecker := &http.Checker{}
	icmpChecker := &icmp.Checker{}
	proxyURL := "http://54.154.45.26:443" 

	config,err := buildNetTestConfig(proxyURL) 
	if err != nil {
		fmt.Printf("Invalid config: %v", err)
		return
	}
	executor := app.NewExecutor(tcpChecker, dnsChecker, httpChecker, icmpChecker, domain.PolicyOptimized)
	result := executor.Run(ctx, config)
	
	var renderer domain.Renderer
	if rsvpConf.textRender{
		renderer = text.NewRenderer()
		if err := renderer.Render(os.Stdout, result); err != nil{
			fmt.Printf("Failed to render: %v", err)
		}
	} else{
		renderer = text.NewTableRenderer()
		if err := renderer.Render(os.Stdout,result); err != nil {
			fmt.Printf("Failed to render: %v", err)
		}
	}
}

func printHeader(){
	fmt.Println("\nRSVP CHECK - Connectivity Diagnostics")
	fmt.Println("══════════════════════════════════════")
}


func buildNetTestConfig(proxyURL string) (domain.NetTestConfig, error) {

	directEndpoints := []domain.Endpoint{
		domain.MustNewICMPEndpoint("1.1.1.1", domain.EndpointTypePublic,"ping 1.1.1.1"),
		domain.MustNewICMPEndpoint("8.8.8.8", domain.EndpointTypePublic,"ping 8.8.8.8"),
		domain.MustNewICMPEndpoint("google.com", domain.EndpointTypePublic,"ping google.com"),
		domain.MustNewDNSEndpoint("insite-eu.gehealthcare.com", domain.EndpointTypePublic,"DNS resolution insite-eu"),
		domain.MustNewDNSEndpoint("insite.gehealthcare.com", domain.EndpointTypePublic,"DNS resolution insite-eu"),
		domain.MustNewDNSEndpoint("google.com", domain.EndpointTypePublic,"DNS resolution google.com"),
		domain.MustNewDNSEndpoint("cloudflare.com", domain.EndpointTypePublic,"DNS resolution claudflare.com"),
		domain.MustNewTCPEndpoint("google.com:443", domain.EndpointTypePublic, "Google HTTPS"),
		domain.MustNewHTTPEndpoint("https://insite-eu.gehealthcare.com:443", 
									domain.EndpointTypePublic, 
									false,
									"",
									"GE Healthcare InSite (direct Internet)"),
	}
	proxyEndpoints := []domain.Endpoint{
		//domain.MustNewICMPEndpoint("54.154.45.26", domain.EndpointTypePublic,"Ping Internet proxy"),
		domain.MustNewHTTPEndpoint("https://insite-eu.gehealthcare.com:443", 
									domain.EndpointTypePublic, 
									true,
									proxyURL,
									"GE Healthcare InSite (via proxy)"),
	}
	
	vpnEndpoints := []domain.Endpoint{
		domain.MustNewICMPEndpoint("150.2.101.89", domain.EndpointTypeVPN,"ping VPN"),
		domain.MustNewICMPEndpoint("82.136.152.65", domain.EndpointTypeVPN,"ping SJUNET"),
		domain.MustNewTCPEndpoint("150.2.101.89:443", domain.EndpointTypeVPN, "VPN endpoint 1"),
		domain.MustNewTCPEndpoint("82.136.152.65:443", domain.EndpointTypeVPN, "VPN endpoint 2"),
	}
	return domain.NewNetTestConfig(
		vpnEndpoints,
		directEndpoints,
		proxyEndpoints,
		proxyURL,       // e.g. "http://proxy.corp:8080"
	)
}