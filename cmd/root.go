package cmd

func Execute() {
	LoadConfig(configFile)
	/* let's go! */
	if helmCmd != "" {
		ProcessHelm()
	}
	if kubectlCmd != "" {
		ProcessKubectl()
	}
}
