package scan

import "fmt"

type hndl struct{}

func (*hndl) Handle(*Event) { fmt.Println("from Handle") }

func ExampleScan() {
	fs := NewOsFs()

	Scan("/somepath", fs, &hndl{})
	// Output:
}
