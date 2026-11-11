// Package pathutil compares repository paths across supported platforms.
package pathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Same(a, b string) bool {
	a, b = filepath.Clean(a), filepath.Clean(b)
	if a == b || runtime.GOOS == "windows" && strings.EqualFold(a, b) {
		return true
	}
	// Git resolves symlinks, including /var -> /private/var on macOS and
	// Windows short paths. Compare filesystem identity before rejecting aliases.
	first, errA := os.Stat(a)
	second, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(first, second)
}
