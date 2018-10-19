package main

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func loadYaml(fileName string) map[string]interface{} {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	v := make(map[string]interface{})
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		log.Fatalln(err)
	}

	return v
}

func mergeYamls(srcFile string, dstFile string) string {
	src := loadYaml(srcFile)
	dest := loadYaml(dstFile)

	mergo.Merge(&dest, src)
	d, err := yaml.Marshal(dest)
	if err != nil {
		panic(err)
	}
	errt := ioutil.WriteFile(dstFile, d, 0644)
	if errt != nil {
		log.Fatalf("Cannot write file: %v\n", errt)
	}

	return dstFile
}

func main() {
	mergeYamls("/home/srua/Code/retail-aws/aws-sandbox-kubernetes/masters.conf.d/ansible.yaml", "/tmp/nodes.yaml")
	mergeYamls("/home/srua/Code/retail-aws/aws-sandbox-kubernetes/masters.conf.d/sizing.yaml", "/tmp/nodes.yaml")
}
