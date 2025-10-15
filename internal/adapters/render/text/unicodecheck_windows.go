//go:build windows

package text

import "syscall"

func isUnicodeSupported() bool {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	getCP := k32.NewProc("GetConsoleOutputCP")
	r1, _, _ := getCP.Call()
	cp := uint32(r1)
	return cp == 65001 // 65001 = UTF-8
}
