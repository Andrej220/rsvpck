package text

import (
	"fmt"
	"io"
	"strings"
)

type colorFunc = func(a ...any) string

const (
	devider1Unicode     = "═"
	devider2Unicode     = "─"
	devider1NoneUnicode = "="
	devider2NoneUnicode = "-"
	minDevLen           = 30
)

func PrintBlock(w io.Writer, name, text string, conf *RenderConfig) error {

	var devLen int = minDevLen
	if minDevLen < len(name) {
		devLen = len(name)
	}
	fmt.Fprintln(w, name)
	fmt.Fprintln(w, strings.Repeat(conf.Divider2, devLen))
	fmt.Fprintln(w, text)
	fmt.Fprintln(w, strings.Repeat(conf.Divider1, devLen))
	return nil
}

func PrintList[T fmt.Stringer](w io.Writer, name string, data []T, conf *RenderConfig) error {
	var devLen int = minDevLen
	if minDevLen < len(name) {
		devLen = len(name)
	}
	fmt.Fprint(w, name)
	fmt.Fprintln(w, strings.Repeat(conf.Divider2, devLen))
	for _, v := range data {
		fmt.Fprintln(w, v.String())
	}
	fmt.Fprintln(w, strings.Repeat(conf.Divider1, devLen))
	return nil
}

