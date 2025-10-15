package text

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"strings"
)

var devider1 string
var devider2 string
var okSym, failSym string

type colorFunc = func(a ...any) string

var (
	green colorFunc
	red   colorFunc
)

const (
	devider1Unicode     = "═"
	devider2Unicode     = "─"
	devider1NoneUnicode = "="
	devider2NoneUnicode = "-"
	minDevLen           = 30
)

func init() {
	green, red = defineGreenRed(unicodeSupported)
	devider1, devider2 = divider(unicodeSupported)
	okSym, failSym = okFailSym(unicodeSupported)
}

func PrintBlock(w io.Writer, name, text string) error {

	var devLen int = minDevLen
	if minDevLen < len(name) {
		devLen = len(name)
	}
	fmt.Fprintln(w, name)
	fmt.Fprintln(w, strings.Repeat(devider2, devLen))
	fmt.Fprintln(w, text)
	fmt.Fprintln(w, strings.Repeat(devider1, devLen))
	return nil
}

func PrintList[T fmt.Stringer](w io.Writer, name string, data []T) error {
	var devLen int = minDevLen
	if minDevLen < len(name) {
		devLen = len(name)
	}
	fmt.Fprint(w, name)
	fmt.Fprintln(w, strings.Repeat(devider2, devLen))
	for _, v := range data {
		fmt.Fprintln(w, v.String())
	}
	fmt.Fprintln(w, strings.Repeat(devider1, devLen))
	return nil
}

func defineGreenRed(supported bool) (colorFunc, colorFunc) {
	if supported {
		green = color.New(color.FgGreen).SprintFunc()
		red = color.New(color.FgRed).SprintFunc()
	} else {
		// fallback: return text as-is
		green = func(a ...any) string { return fmt.Sprint(a...) }
		red = func(a ...any) string { return fmt.Sprint(a...) }
	}
	return green, red
}

func divider(supported bool) (string, string) {
	if supported {
		return devider1Unicode, devider2Unicode
	}
	return devider1NoneUnicode, devider2NoneUnicode
}

func okFailSym(supported bool) (string, string) {
	if supported {
		return green("✓"), red("✗")
	}
	return "OK", "X"
}
