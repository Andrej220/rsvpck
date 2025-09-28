package main

import (
	"time"
	"fmt"
	"strings"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

const delimeterLength = 50

type NetTestConfig struct {
	SiteID          string
	HostName		string
	TestDate 		time.Time
    Timeout         time.Duration
    CheckProxies     []string
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

	out := fmt.Sprintf("Machine uuid: %s, \nHost name: %s, \nDate: %s \n", 
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

	var statusColored string
	switch r.Status {
	case StatusFail:
		statusColored = red + r.Status.String() + reset
		//msg = red + msg + reset
	case StatusPass:
		statusColored = green + r.Status.String() + reset
		//msg = green + msg + reset
	default:
		statusColored = r.Status.String() 
	}

	return fmt.Sprintf("%s: %s \n\tlatency: %s \n\t%s",
		r.TestName, statusColored, lat, msg)
}

func (r *NetTestResult) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func PrintNetTestResult(data []NetTestResult, config *NetTestConfig) {
	fmt.Println(strings.Repeat("=", delimeterLength))
	fmt.Println(config.String())
	for _, v := range(data){
		fmt.Println(v.String())
		fmt.Println(strings.Repeat("", delimeterLength))
	}
}
