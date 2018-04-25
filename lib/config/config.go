package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Root      string
	Blacklist []string
	Depth     int
	Alias     []string
}

func writeToFile(config Config, path string) {
	file, err := yaml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(path, []byte(file), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Init : initialise the captain config file
func Init(config Config, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := strings.TrimSuffix(path, "/config.yaml")
		os.Mkdir(dir, os.ModePerm)
		writeToFile(config, path)
	}
}

// Set : set an alias for a projecft
func Set(path string, config Config, project string, alias string) {
	found := false
	for index, a := range config.Alias {
		if strings.Split(a, "|")[0] == project {
			config.Alias[index] = project + "|" + alias
			found = true
			break
		}
	}
	if !found {
		config.Alias = append(config.Alias, project+"|"+alias)
	}
	writeToFile(config, path)
}

// ParseYAML : parse the config file
func ParseYAML(path string) Config {
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
