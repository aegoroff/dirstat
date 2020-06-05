package sys

import "runtime"

// RunUnderWindows gets whether code running under Microsoft Windows
func RunUnderWindows() bool {
	return runtime.GOOS == "windows"
}
