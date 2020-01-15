package util

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//
// Command line flags and usage message.
//

func init() {
	flag.Usage = func() {
		desc := `Test a WAF is blocking malicious requests.`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] host [host2 ...]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}
}

var Scheme = flag.String("s", "http", "sheme")
var help = flag.Bool("h", false, "print help")

func Flag() []string {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return flag.Args()
}

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

func ParseYamlFile(filename string) (Yaml, error) {
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

func ParseYamlFiles(dirname string) []Yaml {
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
		yaml, err := ParseYamlFile(dirname + "/" + fi.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: %s\n", err)
			continue
		}
		yamls = append(yamls, yaml)
	}

	return yamls
}
