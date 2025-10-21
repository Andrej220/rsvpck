package domain

import(
	"net/url"
	"fmt"
)

type ProxyConfig struct {
	enabled bool
	url     string 
}

func NewProxyConfig(enabled bool, url string) ProxyConfig {
	return ProxyConfig{
		enabled: enabled,
		url:     url,
	}
}

func (p *ProxyConfig) Set(url string){
	p.enabled = true
	p.url = url
}

func (p ProxyConfig) Enabled() bool { 
	return p.enabled 
}

func (p ProxyConfig) URL() string { 
	return p.url 
}

func (p ProxyConfig) MustUseProxy() bool { 
	return p.enabled && p.url != ""
}

func (p ProxyConfig) String() string {
	if p.enabled {
		return fmt.Sprintf("Proxy enabled: %s", p.url)
	}
	return "Proxy disabled"
}

func (p ProxyConfig) IsValid() bool {
	if !p.enabled {
		return true 
	}
	_, err := url.Parse(p.url)
	return err == nil && p.url != ""
}