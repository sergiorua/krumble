package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// exists returns whether the given file or directory exists or not
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateTempDirectory() string {
	dir, err := ioutil.TempDir("", "krumble")
	if err != nil {
		log.Fatalf("Cannot create temp directory: %v\n", err)
	}
	//defer os.RemoveAll(dir)
	return dir
}

func sliceins(arr []string, pos int, elem string) []string {
	if pos < 0 {
		pos = 0
	} else if pos >= len(arr) {
		pos = len(arr)
	}
	out := make([]string, len(arr)+1)
	copy(out[:pos], arr[:pos])
	out[pos] = elem
	copy(out[pos+1:], arr[pos:])
	return out
}

func buildDockerCommand(command string, args ...string) []string {
	tmpdir := fmt.Sprintf("%s:%s", tempDir, tempDir)
	doc := []string{"run", "--rm", "-v", tmpdir, dockerImg, path.Base(command)}

	for i := len(doc) - 1; i >= 0; i-- {
		args = sliceins(args, 0, doc[i])
	}
	return args
}

func runCommandDry(command string, args ...string) error {
	log.Printf("CMD: %s %v\n", command, args)
	return nil
}

func runCommand(command string, args ...string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

	if dockerImg != "" {
		log.Println("Using docker container for tools")
		args = buildDockerCommand(command, args...)
		command = dockerCmd
	}

	if dryrun {
		return runCommandDry(command, args...)
	}
	if debug {
		runCommandDry(command, args...)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	log.Println("Result: " + out.String())

	return err
}
