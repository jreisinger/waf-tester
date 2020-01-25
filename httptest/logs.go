package httptest

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

var logs []LogLine

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

	// Parse log file only once.
	if len(logs) == 0 {
		logs, _ = GetLogLines(logspath)
	}

	for _, l := range logs {
		if l.Transaction.Request.Headers.WafTesterID == t.ID {
			//fmt.Println(l)
			// Print a report.
			t.Logs = append(t.Logs, l)
		}
	}
}

// GetLogLines gets lines of logs from a log file.
func GetLogLines(logspath string) ([]LogLine, error) {
	var logLines []LogLine

	f, err := os.Open(logspath)
	if err != nil {
		return logLines, err
	}
	defer f.Close()

	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		if strings.HasPrefix(line, "{") {
			logLine := parseJSON([]byte(line))
			logLines = append(logLines, logLine)
		}
	}

	return logLines, nil
}

// parseJSON parses JSON ModSecurity audit logs and  returns LogLine.
func parseJSON(data []byte) LogLine {
	var logline LogLine

	// Unmarshal (decode) JSON data.
	if err := json.Unmarshal(data, &logline); err != nil {
		log.Fatalf("JSON decoding failed: %s", err)
	}

	return logline
}
