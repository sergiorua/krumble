package cmd

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
)

var versionVar bool
var configFile string
var kubeconfig string
var dryrun bool
var debug bool
var runOnly string
var dockerImg string
var tempDir string

func init() {
	home := homedir.HomeDir()

	flag.BoolVar(&versionVar, "version", false, "Show version")
	flag.BoolVar(&dryrun, "dry-run", false, "Dry Run")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&configFile, "config", filepath.Join(home, ".krumble.yaml"), "Absolute path to the config file")
	flag.StringVar(&runOnly, "run-only", "all", "Mostly for debug, run only kops, nodes, kubectl or helm")
	flag.StringVar(&dockerImg, "docker", "", "Docker image to use where you have kops, kubectl and helm installed")

	if home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()
	flag.CommandLine.Parse([]string{})
	flag.Set("logtostderr", "true")

	tempDir = CreateTempDirectory()
	log.Printf("Work dir set to %s\n", tempDir)

	if versionVar {
		fmt.Println("Version 0.0.1")
		os.Exit(0)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Sorry, I can't find the config file %s\n", configFile)
		os.Exit(1)
	}
}
