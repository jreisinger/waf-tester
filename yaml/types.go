package yaml

// Some FTW YAML fields (like output.status) can be both an int or array of ints. Adapted from https://github.com/go-yaml/yaml/issues/100#issuecomment-324964723.

type intArray []int

func (a *intArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
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

type stringArray []string

func (a *stringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
					Data    stringArray
				}
				Output struct {
					// can be both int or array of ints
					Status      intArray
					LogContains string `yaml:"log_contains"`
				}
			}
		}
	}
}
