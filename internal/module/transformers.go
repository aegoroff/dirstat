package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
)

func percentTransformer(val interface{}) string {
	v := val.(float64)
	if v >= 90.0 {
		return color.Sprintf("<red>%.2f%%</>", v)
	}
	if v >= 70.0 {
		return color.Sprintf("<yellow>%.2f%%</>", v)
	}
	return fmt.Sprintf("%.2f%%", v)
}

func sizeTransformer(val interface{}) string {
	sz := val.(uint64)
	return humanize.IBytes(sz)
}
