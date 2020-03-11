package yaml

import (
	"testing"
)

func TestIsYaml(t *testing.T) {
	type testpair struct {
		filePath string
		result   bool
	}

	testpairs := []testpair{
		{"file.yaml", true},
		{"file.yml", true},
		{"file.YAML", true},
		{"file.YML", true},
		{"path/to/file.yaml", true},
		{"file.txt", false},
		{"path/to/file.txt", false},
		{"", false},
		{"/", false},
		{".", false},
		{"yaml", false},
	}

	for _, pair := range testpairs {
		result := isYaml(pair.filePath)
		if result != pair.result {
			t.Error(
				"For", pair.filePath,
				"expected", pair.result,
				"got", result,
			)
		}
	}
}
