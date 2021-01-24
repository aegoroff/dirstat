package module

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

const (
	_ int64 = 1 << (10 * iota)
	kbyte
	mbyte
	gbyte
	tbyte
	pbyte
)

// Range defined integer value range
type Range struct {
	// Min value
	Min int64

	// Max value
	Max int64

	size  int64
	count int64
	files files
}

// NewRange creates new *Range
func NewRange(min int64, max int64) *Range {
	return &Range{Min: min, Max: max, files: make(files, 0)}
}

// Less gets whether value specified less then current value
func (r *Range) Less(y rbtree.Comparable) bool {
	return r.Min < y.(*Range).Min
}

// Equal gets whether value specified equal current value
func (r *Range) Equal(y rbtree.Comparable) bool {
	return r.Min == y.(*Range).Min
}

func (r *Range) String() string {
	return fmt.Sprintf("Between %s and %s", humanSize(r.Min), humanSize(r.Max))
}

// Size gets total size of all files that match the Range
func (r *Range) Size() int64 {
	return r.size
}

// Count gets  the number of files that match the Range
func (r *Range) Count() int64 {
	return r.count
}

// Contains defines whether the number specified within range
func (r *Range) Contains(num int64) bool {
	return num >= r.Min && num <= r.Max
}

func newRanges() rbtree.RbTree {
	rs := []*Range{
		NewRange(0, 100*kbyte),
		NewRange(100*kbyte, mbyte),
		NewRange(mbyte, 10*mbyte),
		NewRange(10*mbyte, 100*mbyte),
		NewRange(100*mbyte, gbyte),
		NewRange(gbyte, 10*gbyte),
		NewRange(10*gbyte, 100*gbyte),
		NewRange(100*gbyte, tbyte),
		NewRange(tbyte, 10*tbyte),
		NewRange(10*tbyte, pbyte),
	}
	t := rbtree.New()
	for _, r := range rs {
		t.Insert(r)
	}
	return t
}
