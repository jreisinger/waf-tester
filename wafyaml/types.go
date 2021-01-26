package wafyaml

import yaml "gopkg.in/yaml.v2"

// Some FTW YAML fields (like output.status) can be both an int or array of
// ints. Adapted from
// https://github.com/go-yaml/yaml/issues/100#issuecomment-324964723.

type IntArray []int

func (a *IntArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []int
	err := unmarshal(&multi)
	if err != nil {
		var single int
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []int{single}
	} else {
		*a = multi
	}
	return nil
}

type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

// Yaml represents a WAF test written in YAML format. Test structure is
// compatible with FTWYaml format -
// https://github.com/coreruleset/ftw/blob/master/docs/YAMLFormat.md
type Yaml struct {
	Tests []Test
}

// Test represents a WAF test.
type Test struct {
	Title  string   `yaml:"test_title"`
	Desc   string   `yaml:",omitempty"`
	File   string   `yaml:",omitempty"`
	Tags   []string `yaml:",omitempty"`
	Stages []StageWrapper
}

// StageWrapper is needed to be compatible with FTWYaml format.
type StageWrapper struct {
	Stage Stage
}

// Stage represents HTTP request and response.
type Stage struct {
	Input  Input
	Output Output
}

// Input represents HTTP request.
type Input struct {
	Headers map[string]string `json:"-"`
	Method  string
	URI     string
	// can be both string or array of strings
	Data StringArray
}

// Output represents HTTP response.
type Output struct {
	// can be both int or array of ints
	Status         IntArray
	LogContains    string `yaml:"log_contains,omitempty"`
	LogContainsNot string `yaml:"no_log_contains,omitempty"`
	ExpectError    bool   `yaml:"expect_error,omitempty"`
}

// String implemets stringer interface. Can be used in fmt.Print functions.
func (y Yaml) String() string {
	out, err := yaml.Marshal(&y)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

// Template prints a YAML template for the WAF tests.
func Template() string {
	return Yaml{
		Tests: []Test{
			{
				Title: "SQLi",
				Desc:  "This test expects HTTP response status 403.",
				Tags:  []string{"sqli", "get"},
				Stages: []StageWrapper{
					{Stage: Stage{
						Input: Input{
							Method:  "GET",
							URI:     "?id=1'%20or%20'1'%20=%20'",
							Headers: map[string]string{"User-Agent": "waf-tester"},
						},
						Output: Output{
							Status: IntArray{403},
						},
					},
					},
				},
			},
			{
				Title: "LFI",
				Desc:  "Logs are expected to contain string 930100. If -logs is not used this test will be skipped.",
				Tags:  []string{"lfi", "post"},
				Stages: []StageWrapper{
					{Stage: Stage{
						Input: Input{
							Method: "POST",
							URI:    "/",
							Headers: map[string]string{
								"User-Agent":   "waf-tester",
								"Content-Type": "application/x-www-form-urlencoded",
							},
							Data: StringArray{"arg=../../../etc/passwd&foo=var"},
						},
						Output: Output{
							LogContains: "930100",
						},
					},
					},
				},
			},
			{
				Title: "RCE",
				Desc:  "If both 'status' and 'log_contains' are defined only status is evaluated.",
				Tags:  []string{"rce", "get"},
				Stages: []StageWrapper{
					{Stage: Stage{
						Input: Input{
							Method: "GET",
							URI:    "/?exec=/bin/bash",
						},
						Output: Output{
							Status:      IntArray{403},
							LogContains: "932160",
						},
					},
					},
				},
			},
		},
	}.String()
}
