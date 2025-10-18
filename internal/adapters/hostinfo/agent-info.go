package hostinfo

import(
	//"strings"
	"os/exec"
	"context"
)

type agentData struct{           
	CRM			string
	EnServer	string
	EnPort		string
	ProxyServer	string
	ProxyPort	string
	AgentStatus	string
	SNumber		string
}

func getEnterpriseServer(ctx context.Context) string {

	cmd := exec.CommandContext(ctx, "/opt/InSite/InSiteAgent/bin/GetEnterpriseServer.py")
	b, err := cmd.CombinedOutput()

	if err == nil {
		return string(b)
	}

	return ""
}
