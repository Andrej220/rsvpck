package text

import (
	"fmt"
	"io"
	"strings"
)

const (
	devider1  = "═"
	devider2  = "─"
	minDevLen = 30
)

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
