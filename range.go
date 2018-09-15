package main

type Range struct {
    Min int64
    Max int64
}

func (r Range) contains(f float64) bool {
    return f >= float64(r.Min) && f <= float64(r.Max)
}
