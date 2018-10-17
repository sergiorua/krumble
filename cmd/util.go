package cmd

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
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

func runCommandDry(command string, args ...string) error {
	log.Printf("CMD: %s %v\n", command, args)
	return nil
}

func runCommand(command string, args ...string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer

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
