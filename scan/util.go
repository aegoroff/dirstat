package scan

import (
	"io"
	"log"
)

// Close wraps io.Closer Close func with error handling
func Close(c io.Closer) {
	if c == nil {
		return
	}
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
