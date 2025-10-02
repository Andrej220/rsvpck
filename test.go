package main

import (
	"time"
	"fmt"
	"os"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type NetTestConfig struct {
	SiteID          string
	SN				string
	HostName		string
	TestDate 		time.Time
    Timeout         time.Duration
    CheckProxies    []string
	CheckEndpoints  []string
}

type NetTestResult struct {
    TestName    	string
	TestShortName 	string
    Status      	testStatus  
    Details     	string
    Latency     	time.Duration
    Error       	error
	ErrMsg			string
}

func (c *NetTestConfig) String() string{

	out := fmt.Sprintf("RSvP Connectivity Diagnostic - Site uuid: %s, \nHost name: %s, \nDate: %s ", 
						c.SiteID, c.HostName, c.TestDate.Format("2006-01-02 15:04:05"))
	return out
}

func (r *NetTestResult)String() string{
	lat := latencyToString(r.Latency)
	if r.Latency == 0 {
		lat = "-"
	}

	msg := r.Details
	if msg == "" && r.ErrMsg != "" {
		msg = r.ErrMsg
	}

	return fmt.Sprintf("%s: \t\t%s latency: %s %s",
		r.TestName, r.Status.String(), lat, msg)
}

func (r *NetTestResult) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func PrintNetTestResult(name string, results []NetTestResult, cfg NetTestConfig) {

	table := tablewriter.NewTable(os.Stdout,
        tablewriter.WithAlignment([]tw.Align{tw.AlignLeft, tw.AlignLeft, tw.AlignRight, tw.AlignLeft}),  
        tablewriter.WithRowAutoWrap(tw.WrapNormal),    
        tablewriter.WithHeaderAutoWrap(tw.WrapTruncate), 
        tablewriter.WithMaxWidth(400),                
    )

    table.Header(name, "Status", "Latency", "Details")

    green := color.New(color.FgGreen).SprintFunc()
    red   := color.New(color.FgRed).SprintFunc()
    yellow:= color.New(color.FgYellow).SprintFunc()
    gray  := color.New(color.FgHiBlack).SprintFunc()  

    passCount, failCount, otherCount := 0, 0, 0

    for _, res := range results {
        statusText := res.Status.String()  
        switch res.Status {
        case StatusPass:
             statusText = green("✓ ") + green(res.Status)   
             passCount++
        case StatusFail:
             statusText = red("✗ ") + red(res.Status)     
             failCount++
        case StatusSkipped:
             statusText = gray("⦿ ") + gray(res.Status)    
             otherCount++
        case StatusWarning :
             statusText = yellow("⚠ ") + yellow(res.Status)  
             otherCount++
        default:
             statusText = gray(res.Status)    
             otherCount++
        }

        latStr := latencyToString(res.Latency)
        table.Append(res.TestShortName, statusText, latStr, res.Details)
    }

    table.Configure(func(cfg *tablewriter.Config) {
        cfg.Footer.Alignment.Global = tw.AlignLeft
    })

    //summaryFooter := []any{
    //    "Summary", 
    //    fmt.Sprintf("PASS: %d", passCount), 
    //    fmt.Sprintf("FAIL: %d", failCount), 
    //}
    //table.Footer(summaryFooter...)

    if err := table.Render(); err != nil {
        fmt.Fprintf(os.Stderr, "Error rendering table: %v\n", err)
    }
}

func printText(text string, header string){
	table := tablewriter.NewTable(os.Stdout,
        tablewriter.WithAlignment([]tw.Align{tw.AlignLeft, tw.AlignLeft, tw.AlignRight, tw.AlignLeft}),  
        tablewriter.WithRowAutoWrap(tw.WrapNormal),    
        tablewriter.WithHeaderAutoWrap(tw.WrapTruncate), 
        tablewriter.WithMaxWidth(400),                
    )
    table.Header(header)
    table.Append(text)

    table.Configure(func(cfg *tablewriter.Config) {
        cfg.Footer.Alignment.Global = tw.AlignLeft
    })

    if err := table.Render(); err != nil {
        fmt.Fprintf(os.Stderr, "Error rendering table: %v\n", err)
    }
}
