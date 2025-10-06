package main

import(
	"os/exec"
	"os"
	"strings"
	//"context"
	//"time"
	//"runtime"
)

func getMachineUUID() string{

	cmd := exec.Command("/opt/InSite/InSiteAgent/bin/AgentStatus")
	b, err := cmd.CombinedOutput()
	if err == nil {
		return string(b)
	}

	if b, err := os.ReadFile("/etc/machine-id"); err == nil {
		s := strings.TrimSpace(string(b))
		if s != "" {
			return s
		}
	}

	if b, err := os.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
		s := strings.TrimSpace(string(b))
		if s != "" {
			return s
		}
	}

	return getHostName()
}

func getHostName() string{
	if b, err := os.Hostname(); err == nil {
		s := strings.TrimSpace(string(b))
		if s != "" {
			return s
		}
	}
	return "unknown"
}

//func getRoutePath() NetTestResult{
//	result := NetTestResult{
//		TestName: "Routing table",
//		TestShortName: "Route",
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	var cmd *exec.Cmd 
//	switch runtime.GOOS {
//	case "windows":
//		cmd = exec.CommandContext(ctx, "cmd", "/C", "route print -4 | findstr 0.0.0.0")
//	case "linux":
//		cmd = exec.CommandContext(ctx, "ip", "r", "show",  "default")
//	default:
//		
//	}
//
//	b, err := cmd.CombinedOutput()
//
//	if err != nil {
//		result.Status = StatusFail
//		result.Details = "Failed to get routing table"
//		result.Error = err
//		return result
//	}
//
//	result.Status = StatusUnknown
//	result.Details = indentMultiline(string(b))
//	return result
//}

