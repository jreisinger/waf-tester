package httptest

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

var (
	format  = "%-4v %-10v %-9v %v\n"
	vformat = "  %-9v %v\n"
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

func (t *Test) PrintVerbose() {
	fmt.Printf(format, t.TestStatus, t.Title, t.Method, t.URL)
	fmt.Printf(vformat, "DESC", t.Desc)
	fmt.Printf(vformat, "FILE", t.File)
	fmt.Printf(vformat, "STATUS", t.Status)
	fmt.Printf(vformat, "CODE", t.StatusCode)
	fmt.Printf(vformat, "EXP_CODES", t.ExpectedStatusCodes)
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

func (t *Test) Print() {
	fmt.Printf(format, t.TestStatus, t.Title, t.Method, t.URL)
}
