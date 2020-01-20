package yaml

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	//"github.com/davecgh/go-spew/spew"
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

	//yamlConfig.Tests[0].File = filename

	for i := range yamlConfig.Tests {
		yamlConfig.Tests[i].File = filename
	}

	//spew.Dump(yamlConfig)

	return yamlConfig, nil
}

// ParseFiles parses all YAML files in a directory.
func ParseFiles(dirname string) []Yaml {
	var yamls []Yaml

	walkFunc := func(path string, fi os.FileInfo, err error) error {
		//filename := path + "/" + fi.Name()
		if !fi.IsDir() {
			yaml, err := ParseFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warn parsing %s: %s\n", path, err)
				return err
			}
			yamls = append(yamls, yaml)
		}
		return nil
	}

	filepath.Walk(dirname, walkFunc)

	//fileInfos, err := dir.Readdir(-1) // -1 means return all entries
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "%s\n", err)
	//	os.Exit(1)
	//}

	return yamls
}
