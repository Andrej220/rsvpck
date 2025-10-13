package domain

import (
	"time"
	//"context"
	"github.com/azargarov/go-utils/autostr"
	//"strings"
	//"reflect"
	//"fmt"
)

type HostInfo struct{
	SID			string `string:"include" display:"System ID"`
	Hostname	string `string:"include"`
	SN			string `string:"include" display:"Serial number"`
	OS			string `string:"include" display:"Operating system"`
	RT 			string `string:"include" display:"Routing table"`
	TLSCert		[]TLSCertificate	
}

type NetInfo struct{
	RoutingTable 	string `string:"include"`
}

type TLSCertificate struct {
    Subject   string		`string:"include"`
    Issuer    string		//`string:"include"`
	NotBefore time.Time		`string:"include"`
    NotAfter  time.Time		`string:"include"`
    Valid     bool			`string:"include"`
}

func (t TLSCertificate) String()string{
	autostrCfg := autostr.Config{Separator: autostr.Ptr("\n"), FieldValueSeparator: autostr.Ptr(" : "), PrettyPrint: true}
	return autostr.String(t, autostrCfg)
}

//type TLSCertificateFetcher interface {
//	// addr - "host:443"; serverName used for SNI/hostname verification.
//	GetCertificates(ctx context.Context, addr, serverName string) ([]TLSCertificate, error)
//}


func NewTLSCertificate() []TLSCertificate{
	return []TLSCertificate{}
}

func NewHostInfo() HostInfo{

	return HostInfo{TLSCert: NewTLSCertificate()}

}
