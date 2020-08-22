package httptest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func foundInLogs(t *Test, id string) bool {
	for _, l := range t.Logs {
		for _, m := range l.Transaction.Messages {
			if m.Details.RuleId == id {
				return true
			}
		}
	}
	return false
}

// AddLogs adds related logs to a Test.
func (t *Test) AddLogs(logspath string) {

	logs, logsError := GetLogLines(logspath)
	if logsError != nil {
		log.Printf("problem getting log lines from %s: %v\n", logspath, logsError)
		return
	}

	for _, l := range logs {
		if l.Transaction.Request.Headers.WafTesterID == t.ID {
			//fmt.Println(l)
			// Print a report.
			t.Logs = append(t.Logs, l)
		}
	}
}

// GetLogLines gets lines of WAF logs from URL or file.
func GetLogLines(logspath string) (logs []LogLine, err error) {
	// logspath is URL
	re := regexp.MustCompile(`^http`)
	if re.MatchString(logspath) {
		logs, err = getLogLinesFromURL(logspath)
		return
	}

	// logspath is file
	logs, err = getLogLinesFromFile(logspath)
	return
}

func getLogLinesFromURL(logspath string) (logs []LogLine, err error) {
	lokiUser := os.Getenv("LOKI_USER")
	lokiPass := os.Getenv("LOKI_PASS")
	if lokiUser == "" || lokiPass == "" {
		err = errors.New("env var LOKI_USER or LOKI_PASS are empty")
		return
	}

	// https://github.com/grafana/loki/blob/master/docs/api.md
	cmd := `curl -s -X GET -G -u "$LOKI_USER:$LOKI_PASS" "` + logspath + `/loki/api/v1/query_range" --data-urlencode "query={instance=~\"waf-.*\"}" | jq -r '.data.result | .[] | .values | .[] | .[]'`
	out, err := exec.Command("/bin/sh", "-c", os.ExpandEnv(cmd)).Output() // NOTE: we are ignoring STDERR and exit code
	if err != nil {
		return
	}

	// Convert byte slice to io.Reader (https://stackoverflow.com/a/29747410)
	r := bytes.NewReader(out)

	logs, err = parseLogs(r)

	return
}

func getLogLinesFromFile(logspath string) (logs []LogLine, err error) {
	f, err := os.Open(logspath)
	if err != nil {
		return
	}
	defer f.Close()

	logs, err = parseLogs(f)

	return
}

func parseLogs(r io.Reader) (logs []LogLine, err error) {
	input := bufio.NewScanner(r)
	for input.Scan() {
		line := input.Text()
		if strings.HasPrefix(line, "{") {
			log := parseJSON([]byte(line))
			logs = append(logs, log)
		}
	}
	if input.Err() != nil {
		err = input.Err()
	}
	return
}

// parseJSON parses JSON ModSecurity audit logs.
func parseJSON(data []byte) LogLine {
	var logline LogLine

	// Unmarshal (decode) JSON data.
	if err := json.Unmarshal(data, &logline); err != nil {
		log.Fatalf("JSON decoding failed: %s", err)
	}

	return logline
}
