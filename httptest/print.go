package httptest

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/fatih/color"
)

var (
	format  = "%-4v\t%-20v %-9v %v\n"
	vformat = "  %-10v %v\n"
)

//
// Template for the report that is printed out.
//

const templ = `    {{.Transaction.TimeStamp}}: {{.Transaction.ClientIP}}:{{.Transaction.ClientPort}} --> {{.Transaction.HostIP}}:{{.Transaction.HostPort}}
      --> {{.Transaction.Request.Method}} {{.Transaction.Request.Headers.Host}} {{.Transaction.Request.Uri}}
      <-- {{.Transaction.Response.HttpCode}}
{{range .Transaction.Messages}}        {{.Details.RuleId}}|{{.Message}}
{{end}}`

var report = template.Must(template.New("report").
	Funcs(template.FuncMap{"baseName": baseName}).
	Parse(templ))

// PrintVerbose prints lot of information about a Test.
func (t *Test) PrintVerbose() {
	if !t.Executed {
		return
	}
	t.Print()
	fmt.Printf(vformat, "DESC", t.Desc)
	fmt.Printf(vformat, "FILE", t.File)
	fmt.Printf(vformat, "STATUS", t.Status)
	fmt.Printf(vformat, "CODE", t.StatusCode)
	fmt.Printf(vformat, "EXP_CODES", t.ExpectedStatusCodes)
	fmt.Printf(vformat, "EXP_LOGS", t.LogContains)
	fmt.Printf(vformat, "EXP_NOLOGS", t.LogContainsNot)
	fmt.Printf(vformat, "EXP_ERR", t.ExpectError)
	fmt.Printf(vformat, "ERROR", t.Err)
	fmt.Printf(vformat, "DATA", t.Data)
	fmt.Printf(vformat, "HEADERS", "")
	for k, v := range t.Headers {
		fmt.Printf("    %s: %v\n", k, v)
	}
	fmt.Printf(vformat, "LOGS", "")
	for _, l := range t.Logs {
		if err := report.Execute(os.Stdout, l); err != nil {
			log.Fatal(err)
		}
	}
}

func setTestStatusColor(testStatus string) string {
	var green = color.New(color.FgGreen).SprintFunc()
	var red = color.New(color.FgRed).SprintFunc()
	var yellow = color.New(color.FgYellow).SprintFunc()

	switch testStatus {
	case "OK":
		testStatus = green(testStatus)
	case "FAIL":
		testStatus = red(testStatus)
	case "ERR":
		testStatus = yellow(testStatus)
	}

	return testStatus
}

// Print prints basic information about a Test.
func (t *Test) Print() {
	if !t.Executed {
		return
	}

	testStatus := setTestStatusColor(t.TestStatus)

	fmt.Printf(format, testStatus, t.Title, t.Method, t.URL)
}
