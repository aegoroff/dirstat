package cmd

// Range defined integer value range
type Range struct {
	// Min value
	Min int64

	// Max value
	Max int64
}

func (r Range) contains(f int64) bool {
	return f >= r.Min && f <= r.Max
}
