package scan

import "fmt"

type hndl struct{}

func (h *hndl) Handle(*ScanEvent) { fmt.Println("from Handle") }

func ExampleScan() {
	fs := NewOsFs()

	Scan("/somepath", fs, &hndl{})
	// Output:
}
