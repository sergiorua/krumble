package cmd

import (
	"bytes"
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
		log.Printf("Merging %s snippets\n", file)
		d := mergeYamls(file, dest)
		if debug {
			log.Printf("%s merged into %s\n", file, d)
		}
	}
	return nil
}

// FIXME: merge this function with runCommand
func getKopsConfig(section string, subsection string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer
	var args []string

	if subsection != "" {
		args = []string{"get", section, subsection, "--name", config.Kops.Name, "--state", config.Kops.State, "-oyaml"}
	} else {
		args = []string{"get", section, "--name", config.Kops.Name, "--state", config.Kops.State, "-oyaml"}
	}

	if debug {
		log.Printf("Running: %s %v\n", kopsCmd, args)
	}

	if dryrun {
		return ""
	}
	cmd := exec.Command(kopsCmd, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Cannot run %v %v: %v %v\n", kopsCmd, args, err, stderr.String())
	}

	// WRITE TO Temp FILE NOW
	tmpfile, errt := ioutil.TempFile(tempDir, "kops.*.yaml")
	if errt != nil {
		log.Fatalf("Cannot create temp file: %v\n", errt)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(out.Bytes()); err != nil {
		log.Fatalf("Cannot write output of 'kops get cluster' to %s: %v\n", tmpfile.Name(), err)
		tmpfile.Close()
		return ""
	}
	if debug {
		log.Printf("%s config saved to %s\n", section, tmpfile.Name())
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
	gconf := getKopsConfig("cluster", "")
	mergeKopsSnippets(sourceDir, gconf)
	log.Printf("Config saved to %s\n", gconf)

	kopsReplaceConfig(gconf)
}

func MergeKopsMasterSnippets(sourceDir string) {
	zones := strings.Split(config.Kops.MasterZones, ",")
	for i := range zones {
		s := fmt.Sprintf("master-%s", zones[i])
		gconf := getKopsConfig("ig", s)
		mergeKopsSnippets(sourceDir, gconf)
		log.Printf("Config saved to %s\n", gconf)
		kopsReplaceConfig(gconf)
	}
}

func MergeKopsNodeSnippets(sourceDir string) {
	gconf := getKopsConfig("ig", "nodes")
	mergeKopsSnippets(sourceDir, gconf)
	kopsReplaceConfig(gconf)
}
