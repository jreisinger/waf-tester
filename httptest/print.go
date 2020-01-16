package httptest

import "fmt"

var format = "%-4s %-9s %s\n"

func (t *Test) PrintVerbose() {
	fmt.Printf(format, t.Status, t.Method, t.URL)
	fmt.Printf("  %s\n", t.Desc)
	for k, v := range t.Headers {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

func (t *Test) Print() {
	fmt.Printf(format, t.Status, t.Method, t.URL)
}
