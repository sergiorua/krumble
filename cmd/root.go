package cmd

var vpc AwsVpc
var subnets []AwsSubnet
var utility_subnets []AwsSubnet

func Execute() {
	LoadConfig(configFile)

	vpc, subnets, utility_subnets = ProcessHooks()

	ProcessKops()

	if helmCmd != "" {
		ProcessHelm()
	}
	if kubectlCmd != "" {
		ProcessKubectl()
	}
}
