package cmd

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

func RunKubectl(entry Kubectl) error {
	args := []string{"apply", "-f", entry.URL}
	err := runCommand(kubectlCmd, args...)
	return err
}

func LoadKubeYamlFromUrl(url string) []byte {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	buf := new(bytes.Buffer)

	if err != nil {
		log.Printf("Could not download from " + url)
		return make([]byte, 0)
	}
	defer response.Body.Close()
	buf.ReadFrom(response.Body)
	respByte := buf.Bytes()

	return respByte
}

func LoadKubeYamlFromFile(fpath string) []byte {
	dat, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Printf("Cannot open file %s: %v\n", fpath, err)
		return make([]byte, 0)
	}

	return dat
}

func ProcessKubectl() {
	for i := range config.Kubectl {
		log.Printf("Processing entry %v\n", config.Kubectl[i].Name)
		err := RunKubectl(config.Kubectl[i])
		if err != nil {
			log.Printf("Cannot apply %s with %s\n", config.Kubectl[i].Name, config.Kubectl[i].URL)
		}
	}
}
