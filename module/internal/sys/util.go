package sys

import (
	"io"
	"log"
	"runtime"
)

// RunUnderWindows gets whether code running under Microsoft Windows
func RunUnderWindows() bool {
	return runtime.GOOS == "windows"
}

// Close wraps io.Closer Close func with error handling
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
