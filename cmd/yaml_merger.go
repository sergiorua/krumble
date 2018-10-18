package cmd

import (
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

func mergeKopsSnippets(dir string, dest string) error {
	files, err := filepath.Glob(path.Join(dir, "*.yaml"))
	if err != nil {
		log.Printf("Cannot read directory %s\n", dir)
		return err
	}
	for _, file := range files {
		log.Printf("Merging %s with cluster\n", file)
		mergeYamls(file, dest)
	}
	return nil
}

func getKopsConfig(section string) string {
	// ie: kops get cluster --name=%s -oyaml
	args := []string{"get", section, "--name", config.Kops.Name, "--state", config.Kops.State, "-oyaml"}
	log.Printf("Running: %s %v\n", kopsCmd, args)
	cmd := exec.Command(kopsCmd, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Cannot run %v %v\n", kopsCmd, args)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	contents, errc := ioutil.ReadAll(stdout)
	if errc != nil {
		log.Fatalf("Reading contents of 'kops get' command: %v\n", errc)
		return ""
	}

	// WRITE TO FILE NOW
	tmpfile, errt := ioutil.TempFile("", "kops.*.yaml")
	if errt != nil {
		log.Fatalf("Cannot create temp file: %v\n", errt)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(contents); err != nil {
		log.Fatalf("Cannot write output of 'kops get cluster' to %s: %v\n", tmpfile.Name(), err)
		tmpfile.Close()
		return ""
	}
	return tmpfile.Name()
}

func kopsReplaceConfig(gconf string) error {
	//kops replace --name $CLUSTER_NAME -f $TMPFILE
	args := []string{"replace", "--name", config.Kops.Name, "--state", config.Kops.State, "-f", gconf}
	err := runCommand(kopsCmd, args...)
	if err != nil {
		log.Fatalf("kops replace failed with error %v\n", err)
		return err
	}
	os.Remove(gconf)
	return nil
}

func MergeKopsClusterSnippets(sourceDir string) {
	gconf := getKopsConfig("cluster")
	mergeKopsSnippets(sourceDir, gconf)
	log.Printf("Config saved to %s\n", gconf)

	kopsReplaceConfig(gconf)
}

func MergeKopsMasterSnippets(sourceDir string) {
	zones := strings.Split(config.Kops.MasterZones, ",")
	for i := range zones {
		s := fmt.Sprintf("master-%s", zones[i])
		gconf := getKopsConfig(s)
		mergeKopsSnippets(sourceDir, gconf)
		log.Printf("Config saved to %s\n", gconf)
	}

	kopsReplaceConfig(gconf)
}

func MergeKopsNodeSnippets(sourceDir string) {
	gconf := getKopsConfig("ig")
	mergeKopsSnippets(sourceDir, gconf)
	log.Printf("Config saved to %s\n", gconf)

	kopsReplaceConfig(gconf)
}
