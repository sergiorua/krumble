package cmd

import (
	"log"
	"os/exec"
)

func runCommandDry(command string, args ...string) error {
	log.Printf("CMD: %s %v\n", command, args)
	return nil
}

func runCommand(command string, args ...string) error {
	if dryrun {
		return runCommandDry(command, args...)
	}
	cmd := exec.Command(command, args...)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", stdout)

	return err
}
