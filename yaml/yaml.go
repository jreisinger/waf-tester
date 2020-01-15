package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Yaml struct {
	Tests []struct {
		Desc   string
		Stages []struct {
			Stage struct {
				Input struct {
					Headers map[string]string `json:"-"`
					Method  string
					URI     string
				}
			}
		}
	}
}

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
		yaml, err := ParseFile(dirname + "/" + fi.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: %s\n", err)
			continue
		}
		yamls = append(yamls, yaml)
	}

	return yamls
}
