package httptest

import "fmt"

var (
	format  = "%-4v %-10v %-9v %v\n"
	vformat = "  %-9v %v\n"
)

func (t *Test) PrintVerbose() {
	fmt.Printf(format, t.TestStatus, t.Title, t.Method, t.URL)
	fmt.Printf(vformat, "TITLE", t.Title)
	fmt.Printf(vformat, "DESC", t.Desc)
	fmt.Printf(vformat, "FILE", t.File)
	fmt.Printf(vformat, "STATUS", t.Status)
	fmt.Printf(vformat, "CODE", t.StatusCode)
	fmt.Printf(vformat, "EXP_CODES", t.ExpectedStatusCodes)
	fmt.Printf(vformat, "ERROR", t.Err)
	fmt.Printf(vformat, "DATA", t.Data)
	fmt.Printf(vformat, "HEADERS", t.Data)
	for k, v := range t.Headers {
		fmt.Printf("    %s: %v\n", k, v)
	}
}

func (t *Test) Print() {
	fmt.Printf(format, t.TestStatus, t.Title, t.Method, t.URL)
}
