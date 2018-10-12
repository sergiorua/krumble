package cmd

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func RunHelm(tmpfile string) error {
	runCommand(helmCmd, "--file", tmpfile, "sync")
	return nil
}

func HelmCleanUp(tmpfile string) {
	os.Remove(tmpfile)
}

func HelmWriteConfig() (string, error) {
	d, err := yaml.Marshal(config.Helm)
	if err != nil {
		log.Fatalf("error: %v", err)
		return "", err
	}

	tmpfile, errt := ioutil.TempFile("", "helm.*.yaml")
	if errt != nil {
		log.Fatalf("error: creating temp file: %v\n", errt)
		return "", errt
	}

	log.Printf("Writing helm config to %s\n", tmpfile.Name())

	defer tmpfile.Close()

	_, err2 := tmpfile.Write(d)
	if err2 != nil {
		log.Fatalf("error: %v", err2)
		return "", err2
	}

	return tmpfile.Name(), nil
}
func ProcessHelm() error {
	hc, err := HelmWriteConfig()
	if err != nil {
		return err
	}

	err = RunHelm(hc)
	if err != nil {
		log.Fatalf("Helmfile could not be process\n")
	}
	HelmCleanUp(hc)
	return err
}
