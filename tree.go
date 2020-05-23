package main

type namedInt64 struct {
	name  string
	value int64
}

type statItem struct {
	name  string
	size  int64
	count int64
}

type namedInts64 []*namedInt64

func (x namedInts64) Len() int {
	return len(x)
}

func (x namedInts64) Less(i, j int) bool {
	return x[i].value < x[j].value
}

func (x namedInts64) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x *statItem) LessThan(y interface{}) bool {
	return x.size < (y.(*statItem)).size
}

func (x *statItem) EqualTo(y interface{}) bool {
	return x.size == (y.(*statItem)).size
}

func (x namedInt64) LessThan(y interface{}) bool {
	return x.value < (y.(namedInt64)).value
}

func (x namedInt64) EqualTo(y interface{}) bool {
	return x.value == (y.(namedInt64)).value
}
