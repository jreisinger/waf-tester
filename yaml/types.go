package yaml

// Yaml represents a WAF test written in YAML format.
type Yaml struct {
	Tests []struct {
		Title  string `yaml:"test_title"`
		Desc   string
		File   string
		Stages []struct {
			Stage struct {
				Input struct {
					Headers map[string]string `json:"-"`
					Method  string
					URI     string
					Data    string
				}
			}
		}
	}
}
