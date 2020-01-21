package httptest

// Test represents an HTTP test. It contains both request and response fields.
type Test struct {
	Title              string
	Desc               string
	File               string
	Method             string
	Scheme             string
	Host               string
	Path               string // URI
	URL                string // scheme + host + Path
	Headers            map[string]string
	Data               string
	Err                error
	StatusCode         int
	ExpectedStatusCode int
	Status             string
	TestStatus         string
}