package yaml

// Yaml represents a WAF test written in YAML format.
type Yaml struct {
	Tests []struct {
		Desc   string
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
