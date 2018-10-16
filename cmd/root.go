package cmd

var vpc AwsVpc
var subnets []AwsSubnet

func Execute() {
	LoadConfig(configFile)

	/* let's go! */

	/* discover vpc and subnets */
	vpc, subnets = ProcessHooks()

	ProcessKops()

	if helmCmd != "" {
		ProcessHelm()
	}
	if kubectlCmd != "" {
		ProcessKubectl()
	}
}
