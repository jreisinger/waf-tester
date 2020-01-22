package httptest

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

// AddLogs adds related logs to a Test.
func (t *Test) AddLogs(logspath string) {
	logs := getLogLines(logspath)
	for _, l := range logs {
		if l.Transaction.Request.Headers.WafTesterID == t.ID {
			//fmt.Println(l)
			// Print a report.
			t.Logs = append(t.Logs, l)
		}
	}
}

func getLogLines(logspath string) []LogLine {
	logs, err := parseLogFile(logspath)
	if err != nil {
		log.Fatal(err)
	}
	return logs
}

func parseLogFile(filename string) ([]LogLine, error) {
	var logLines []LogLine

	f, err := os.Open(filename)
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
