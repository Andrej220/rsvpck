package text

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	//"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type RenderConfig struct {

	ForceASCII bool
	ForceUnicode *bool // nil = auto, true/false = force
	Unicode bool
	Color   bool
	OkSym, FailSym string
	Divider1, Divider2 string
	Green, Red colorFunc
	TableSymbols *tw.SymbolCustom
}

const maxCharPerError = 150

type Option func(*RenderConfig)

func WithForceASCII(v bool) Option        { return func(c *RenderConfig) { c.ForceASCII = v } }
func WithForceUnicode(v bool) Option      { return func(c *RenderConfig) { c.ForceUnicode = &v } }

func NewRenderConfig(opts ...Option) *RenderConfig {
	c := &RenderConfig{}
	for _, opt := range opts { opt(c) }

	autoUnicode := isUnicodeSupported()

	if v := os.Getenv("RSVPCK_UNICODE"); v != "" && c.ForceUnicode == nil {
		vv := strings.TrimSpace(strings.ToLower(v))
		switch vv {
		case "1", "true", "on", "yes": t := true; c.ForceUnicode = &t
		case "0", "false", "off", "no": f := false; c.ForceUnicode = &f
		}
	}

	switch {
	case c.ForceASCII:
		c.Unicode = false
	case c.ForceUnicode != nil:
		c.Unicode = *c.ForceUnicode
	default:
		c.Unicode = autoUnicode
	}

	c.initTheme()
	return c
}


func (c *RenderConfig) initTheme() {
	if c.Unicode {
		color.NoColor = false // enable ANSI
		c.Green = color.New(color.FgGreen).SprintFunc()
		c.Red   = color.New(color.FgRed).SprintFunc()
		c.OkSym, c.FailSym = c.Green("✓"), c.Red("✗")
		c.Divider1, c.Divider2 = "═", "─"
		c.TableSymbols = tw.NewSymbolCustom("Box").
			WithRow("─").WithColumn("│").
			WithTopLeft("┌").WithTopMid("┬").WithTopRight("┐").
			WithMidLeft("├").WithCenter("┼").WithMidRight("┤").
			WithBottomLeft("└").WithBottomMid("┴").WithBottomRight("┘")
	} else {
		color.NoColor = true // disable ANSI
		c.OkSym, c.FailSym = "OK", "X"
		c.Divider1, c.Divider2 = "=", "-"
		c.Green = func(a ...any) string { return fmt.Sprint(a...) }
		c.Red   = func(a ...any) string { return fmt.Sprint(a...) }
		c.TableSymbols = tw.NewSymbolCustom("ASCII").
			WithRow("-").WithColumn("|").
			WithTopLeft("+").WithTopMid("+").WithTopRight("+").
			WithMidLeft("+").WithCenter("+").WithMidRight("+").
			WithBottomLeft("+").WithBottomMid("+").WithBottomRight("+")
	}
}