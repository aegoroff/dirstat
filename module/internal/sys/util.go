package sys

import (
	"io"
	"log"
)

// Close wraps io.Closer Close func with error handling
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
