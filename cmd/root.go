package cmd

var vpc AwsVpc
var subnets []AwsSubnet
var utility_subnets []AwsSubnet

func Execute() {
	LoadConfig(configFile)

	/* let's go! */

	/* discover vpc and subnets */
	vpc, subnets, utility_subnets = ProcessHooks()

	ProcessKops()

	if helmCmd != "" {
		ProcessHelm()
	}
	if kubectlCmd != "" {
		ProcessKubectl()
	}
}
