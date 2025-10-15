//go:build !windows

package text

import (
	"os"
	"strings"
)

func isUnicodeSupported() bool {
	for _, key := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
		if val := os.Getenv(key); strings.Contains(strings.ToUpper(val), "UTF-8") {
			return true
		}
	}
	return false
}
