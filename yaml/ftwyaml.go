package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// ParseFile parses a single YAML file.
func ParseFile(filename string) (Yaml, error) {
	var yamlConfig Yaml

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return yamlConfig, err
	}

	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return yamlConfig, err
	}

	return yamlConfig, nil
}

// ParseFiles parses all YAML files in a directory.
func ParseFiles(dirname string) []Yaml {
	dir, err := os.Open(dirname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer dir.Close()

	var yamls []Yaml

	fileInfos, err := dir.Readdir(-1) // -1 means return all entries
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	for _, fi := range fileInfos {
		filename := dirname + "/" + fi.Name()
		yaml, err := ParseFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn parsing %s: %s\n", filename, err)
			continue
		}
		yamls = append(yamls, yaml)
	}

	return yamls
}
