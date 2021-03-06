package cmd

var vpc AwsVpc
var subnets []AwsSubnet
var utility_subnets []AwsSubnet

func Execute() {
	LoadConfig(configFile)

	vpc, subnets, utility_subnets = ProcessHooks()

	if runOnly == "all" || runOnly == "pre_exec" {
		ProcessExec(config.PreExec)
	}

	if runOnly == "all" || runOnly == "kops" {
		ProcessKops()
	}

	if runOnly == "nodes" {
		KopsNodesUp()
	}

	if runOnly == "all" || runOnly == "helm" {
		if helmfileCmd != "" {
			ProcessHelm()
		}
	}

	if runOnly == "all" || runOnly == "exec" {
		ProcessExec(config.Exec)
	}

	if runOnly == "all" || runOnly == "kubectl" {
		if kubectlCmd != "" {
			ProcessKubectl()
		}
	}

	if runOnly == "all" || runOnly == "post_exec" {
		ProcessExec(config.PostExec)
	}

}
