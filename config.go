package main

import (
	"os"
	"log"
	"regexp"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type githubConfig struct {
	Owner string   `yaml:"owner"`
	Repo string    `yaml:"repo"`
	Regexp string  `yaml:"regexp"`
}

type folderConfig struct {
	Name string    `yaml:"name"`
	Info string    `yaml:"info"`
	Path string    `yaml:"path"`
	Regexp string  `yaml:"regexp"`
}

type regexpConfig struct {
	Name string    `yaml:"name"`
	Info string    `yaml:"info"`
	Path string    `yaml:"path"`
	Regexp string  `yaml:"regexp"`
	Date string    `yaml:"date"`
	Format string  `yaml:"format"`
}

type rootConfig struct {
	Interval int           `yaml:"update_interval"`
	Path string            `yaml:"metrics_path"`
	Github []githubConfig  `yaml:"github"`
	Folder []folderConfig  `yaml:"folder"`
	Regexp []regexpConfig  `yaml:"regexp"`
}

var config rootConfig

func validateRegexp(section, pattern string) {
	_,err := regexp.Compile(pattern)
	if err != nil {
		log.Print(section, ": ", err)
		os.Exit(1)
	}
}

func readConfigFile(fileName string) {
	var section string

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Error reading configuration file: %s\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Printf("Error parsing configuration file: %s\n", err)
		os.Exit(1)
	}

	for _,v := range config.Github {
		section = "github." + v.Repo + ".regexp"
		validateRegexp(section, v.Regexp)
	}

	for _,v := range config.Folder {
		section = "folder." + v.Name + ".regexp"
		validateRegexp(section, v.Regexp)
	}

	for _,v := range config.Regexp {
		section = "regexp." + v.Name + ".regexp"
		validateRegexp(section, v.Regexp)
		section = "regexp." + v.Name + ".date"
		validateRegexp(section, v.Date)
	}
}
