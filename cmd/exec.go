package cmd

import (
	"log"
	"os"
	"os/user"
	"strings"
)

func userHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func setEnvironment(cmd Exec) {
	log.Println("Setting up environment variables: ")
	for i := range cmd.Env {
		e := cmd.Env[i]
		log.Printf("\t%s=%s\n", e.Name, e.Value)
		os.Setenv(e.Name, e.Value)
	}
	// Set home dir
	os.Setenv("HOME", userHomeDir())
}

func execCommand(cmd Exec) error {
	currDir, err := os.Getwd()
	if err != nil {
		log.Printf("os.Getwd() failed: %v\n", err)
		return err
	}
	if (debug || dryrun) && cmd.Rundir != "" {
		log.Printf("Changing to directory %s\n", cmd.Rundir)
	}

	if cmd.Rundir != "" && !dryrun {
		err := os.Chdir(cmd.Rundir)
		if err != nil {
			log.Printf("Cannot change to directory %s: %v\n", cmd.Rundir, err)
			return err
		}
	}

	log.Printf("Running %v\n", cmd.Command)
	fullCmd := strings.Fields(cmd.Command)

	command := fullCmd[0]
	args := append(fullCmd[:0], fullCmd[1:]...)

	err = runCommand(command, args...)
	if err != nil {
		log.Printf("Exec '%s' failed: %v\n", fullCmd, err)
	}
	if debug {
		log.Printf("Changing to directory %s\n", currDir)
	}
	os.Chdir(currDir)
	return err
}

func ProcessExec() error {
	for i := range config.Exec {
		if debug {
			log.Printf("Processing command %v\n", config.Exec[i])
		}
		setEnvironment(config.Exec[i])
		execCommand(config.Exec[i])
	}

	return nil
}
