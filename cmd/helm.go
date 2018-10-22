package cmd

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func RunHelm(tmpfile string) error {
	runCommand(helmfileCmd, "--file", tmpfile, "sync")
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

	tmpfile, errt := ioutil.TempFile(tempDir, "helm.*.yaml")
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

// FIXME: check if service account already exists
func installHelm() error {
	var args []string
	var err error

	/* Don't re-install */
	if getPodStatus("tiller-deploy", "kube-system") != "Unknown" {
		log.Println("Tiller already installed")
		return nil
	}

	account := getServiceAccount("tiller", "kube-system")
	if account.Name == "" {
		args = []string{"create", "serviceaccount", "--namespace", "kube-system", "tiller"}
		err = runCommand(kubectlCmd, args...)
		if err != nil {
			log.Printf("Aborting tiller installation: %v\n", err)
			return err
		}
	}

	args = []string{"create", "clusterrolebinding", "tiller-cluster-rule", "--clusterrole=cluster-admin", "--serviceaccount=kube-system:tiller"}
	err = runCommand(kubectlCmd, args...)
	if err != nil {
		log.Printf("Aborting tiller installation: %v\n", err)
		return err
	}
	args = []string{"init", "--service-account", "tiller"}
	err = runCommand(helmCmd, args...)
	if err != nil {
		log.Printf("Aborting tiller installation: %v\n", err)
		return err
	}
	log.Printf("Waiting for helm/tiller...")
	time.Sleep(30 * time.Second)
	return nil
}

func ProcessHelm() error {
	ierr := installHelm()
	if ierr != nil {
		return ierr
	}

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
