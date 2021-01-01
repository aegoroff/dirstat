package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

func percentTransformer(val interface{}) string {
	return fmt.Sprintf("%.2f%%", val)
}

func sizeTransformer(val interface{}) string {
	sz := val.(uint64)
	return humanize.IBytes(sz)
}
