package hostinfo

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"github.com/azargarov/rsvpck/internal/domain"
)
var windowsRoutingTableCommand = []string{"cmd", "/C", "route print -4 | findstr 0.0.0.0"}
var linuxRoutingTableCommand = []string{"ip", "r", "show",  "default"}


func GetCRMInfo(ctx context.Context) domain.HostInfo{

	info := domain.NewHostInfo()
	info.Hostname = getHostname()
	info.OS = runtime.GOOS
	info.SID = getSID(ctx)
	info.RT = string(getRoutingTable())

	return info
}

func getSID(ctx context.Context,) string{

	cmd := exec.CommandContext(ctx,"/opt/InSite/InSiteAgent/bin/AgentStatus")
	b, err := cmd.CombinedOutput()
	if err == nil {
		return string(b)
	}

	hostname :=getHostname()
	if hostname != ""{
		return hostname
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
	return "unknown"
}

func getHostname() string{
	if b, err := os.Hostname(); err == nil {
		s := strings.TrimSpace(string(b))
		if s != "" {
			return s
		}
	}
	return "unknown"
}

func getRoutingTable() []byte {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cmd *exec.Cmd 
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, windowsRoutingTableCommand[0],windowsRoutingTableCommand[1:]...)//"cmd", "/C", "route print -4 | findstr 0.0.0.0")
	case "linux":
		cmd = exec.CommandContext(ctx, linuxRoutingTableCommand[0], linuxRoutingTableCommand[1:]...)//"ip", "r", "show",  "default")
	default:
		return []byte{}
	}

	b, err := cmd.CombinedOutput()

	if err != nil {
		return []byte{}
	}
	return b
}


