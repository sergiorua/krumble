package cmd

import (
	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Global struct {
	ClusterName string `yaml:"cluster_name"`
	Domain      string `yaml:"domain"`
	Environment string `yaml:"environment"`
	Provider    string `yaml:"provider"`
	Aws         struct {
		Region string `yaml:"region"`
		VpcID  struct {
			Hook string `yaml:"hook"`
		} `yaml:"vpc_id"`
		Subnets struct {
			Filters struct {
				TagName string `yaml:"tag:Name"`
			} `yaml:"filters"`
			Hook string `yaml:"hook"`
		} `yaml:"subnets"`
	} `yaml:"aws"`
}

type Kubectl struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	Namespace string `yaml:"namespace"`
}

type ConfigData struct {
	Global  Global
	Kubectl []Kubectl
	Helm    interface{}
}

/* global holding all the yaml config */
var config ConfigData
var helmCmd string
var kubectlCmd string

func CmdLookPath(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal("Could not find %s in the path\n", cmd)
		return ""
	}

	log.Printf("using %s\n", path)
	return path
}

func LoadConfig(configFile string) {
	helmCmd = CmdLookPath("helmfile")
	kubectlCmd = CmdLookPath("kubectl")
	log.Printf("Loading config file from %s\n", configFile)

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error opening %s: #%v\n", configFile, err)
		os.Exit(1)
	}

	c := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	err = mapstructure.Decode(c["global"], &config.Global)
	if err != nil {
		log.Fatalf("Error decoding global section: %v\n", err)
	}

	err = mapstructure.Decode(c["kubectl"], &config.Kubectl)
	if err != nil {
		log.Fatalf("Error decoding global section: %v\n", err)
	}
	config.Helm = c["helm"]
}
