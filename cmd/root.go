package cmd

import (
	"fmt"
)

func Execute() {
	fmt.Printf("Calling config.LoadConfig\n")
	LoadConfig(configFile)

	/* let's go! */
	ProcessKubectl()
}
