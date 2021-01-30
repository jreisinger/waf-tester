package waftest

import "path/filepath"

// Test represents an HTTP test. It contains both request and response fields + additional fields.
type Test struct {
	ID                  string
	Title               string
	Desc                string
	Tags                []string
	File                string
	Method              string
	Scheme              string
	Host                string
	Path                string // URI
	URL                 string // scheme + host + Path
	Headers             map[string]string
	Data                []string
	Err                 error
	StatusCode          int // ex: 403
	ExpectedStatusCodes []int
	Status              string // ex: 403 Forbidden
	TestStatus          string
	Logs                []LogLine
	LogContains         string
	LogContainsNot      string
	ExpectError         bool
	Executed            bool
}

//
// JSON fields from the ModSecurity audit log that we are interested in.
//

// LogLine is a line from logs.
type LogLine struct {
	Transaction Transaction
}

// Transaction is transaction field in logs.
type Transaction struct {
	TimeStamp  string `json:"time_stamp"`
	ClientIP   string `json:"client_ip"`
	ClientPort int    `json:"client_port"`
	HostIP     string `json:"host_ip"`
	HostPort   int    `json:"host_port"`
	Messages   []Message
	Request    Request
	Response   Response
	Producer   Producer
}

// Producer is producer field in logs.
type Producer struct {
	Modsecurity string
	Connector   string
	Components  []string
}

// Request is request field in logs.
type Request struct {
	Method  string
	URI     string
	Headers Headers
}

// Headers is headers field in logs.
type Headers struct {
	Host        string
	WafTesterID string `json:"waf-tester-id"`
}

// Response is response field in logs.
type Response struct {
	HTTPCode int `json:"http_code"`
}

// Message is message field in logs.
type Message struct {
	Message string
	Details Details
}

// Details is details field in logs.
type Details struct {
	RuleID     string
	File       string
	LineNumber string
	Data       string
	Match      string
	Tags       []string
}

func baseName(path string) string {
	return filepath.Base(path)
}
