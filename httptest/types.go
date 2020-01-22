package httptest

import "path/filepath"

// Test represents an HTTP test. It contains both request and response fields.
type Test struct {
	ID                  string
	Title               string
	Desc                string
	File                string
	Method              string
	Scheme              string
	Host                string
	Path                string // URI
	URL                 string // scheme + host + Path
	Headers             map[string]string
	Data                []string
	Err                 error
	StatusCode          int
	ExpectedStatusCodes []int
	Status              string
	TestStatus          string
	Logs                []LogLine
}

//
// JSON fields from the ModSecurity audit log that we are interested in.
//

type LogLine struct {
	Transaction Transaction
}

type Transaction struct {
	TimeStamp  string `json:"time_stamp"`
	ClientIP   string `json:"client_ip"`
	ClientPort int    `json:"client_port"`
	HostIP     string `json:"host_ip"`
	HostPort   int    `json:"host_port"`
	Messages   []Message
	Request    Request
	Response   Response
}

type Request struct {
	Method  string
	Uri     string
	Headers Headers
}

type Headers struct {
	Host        string
	WafTesterID string `json:"waf-tester-id"`
}

type Response struct {
	HttpCode int `json:"http_code"`
}

type Message struct {
	Message string
	Details Details
}

type Details struct {
	RuleId     string
	File       string
	LineNumber string
	Data       string
}

func baseName(path string) string {
	return filepath.Base(path)
}
