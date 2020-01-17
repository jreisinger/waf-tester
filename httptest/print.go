package httptest

import "fmt"

var (
	format  = "%-4s %-9s %s\n"
	vformat = "  %-6s: %v\n"
)

func (t *Test) PrintVerbose() {
	fmt.Printf(format, t.TestStatus, t.Method, t.URL)
	fmt.Printf(vformat, "DESC", t.Desc)
	fmt.Printf(vformat, "STATUS", t.Status)
	fmt.Printf(vformat, "ERROR", t.Err)
	for k, v := range t.Headers {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

func (t *Test) Print() {
	fmt.Printf(format, t.TestStatus, t.Method, t.URL)
}
