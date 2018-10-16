package cmd

var vpc AwsVpc
var subnets []AwsSubnet

func Execute() {
	LoadConfig(configFile)

	getK8sNodes()
	os.Exit(0)
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
