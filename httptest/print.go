package httptest

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/fatih/color"
)

var (
	format  = "%-4v %-10v %-9v %v\n"
	vformat = "  %-10v %v\n"
)

//
// Template for the report that is printed out.
//

const templ = `    {{.Transaction.TimeStamp}}: {{.Transaction.ClientIP}}:{{.Transaction.ClientPort}} --> {{.Transaction.HostIP}}:{{.Transaction.HostPort}}
      --> {{.Transaction.Request.Method}} {{.Transaction.Request.Uri}} {{.Transaction.Request.Headers.Host}}
      <-- {{.Transaction.Response.HttpCode}}
{{range .Transaction.Messages}}        {{.Details.RuleId}}|{{.Message}}
{{end}}`

var report = template.Must(template.New("report").
	Funcs(template.FuncMap{"baseName": baseName}).
	Parse(templ))

// PrintVerbose prints lot of information about a Test.
func (t *Test) PrintVerbose() {
	t.Print()
	fmt.Printf(vformat, "DESC", t.Desc)
	fmt.Printf(vformat, "FILE", t.File)
	fmt.Printf(vformat, "STATUS", t.Status)
	fmt.Printf(vformat, "CODE", t.StatusCode)
	fmt.Printf(vformat, "EXP_CODES", t.ExpectedStatusCodes)
	fmt.Printf(vformat, "EXP_LOGS", t.LogContains)
	fmt.Printf(vformat, "EXP_NOLOGS", t.LogContainsNot)
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

// Print prints basic information about a Test.
func (t *Test) Print() {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	testStatus := t.TestStatus

	switch t.TestStatus {
	case "OK": testStatus = green(testStatus)
	case "FAIL": testStatus = red(testStatus)
	case "ERR": testStatus = yellow(testStatus)
	}

	fmt.Printf(format, testStatus, t.Title, t.Method, t.URL)
}
