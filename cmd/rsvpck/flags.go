package main

import (
	"flag"
)

type rsvpckConf struct {
	textRender  bool
	tableRender bool
	forseASCII bool
}

func NewRsvpckConf() rsvpckConf {
	return rsvpckConf{tableRender: true}
}

func (r *rsvpckConf) setTextRenderOn() {
	r.tableRender = false
	r.textRender = true
}

func (r *rsvpckConf) setTextRenderOff() {
	r.tableRender = true
	r.textRender = false
}

func (r *rsvpckConf) SetRender(textRender bool) {
	if textRender {
		r.setTextRenderOn()
		return
	}
	r.setTextRenderOff()
}

func parseFlagsToConfig() *rsvpckConf {
	txtRender := flag.Bool("text", false, "render connectivity info as text. Default table")
	flagForceASCII := flag.Bool("ascii", false, "Force ASCII-only output (no Unicode symbols)")
	flag.Parse()

	r := NewRsvpckConf()
	r.SetRender(*txtRender)
	r.forseASCII = *flagForceASCII
	return &r
}
