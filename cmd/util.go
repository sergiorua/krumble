package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

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
