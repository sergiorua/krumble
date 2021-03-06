package cmd

import (
	"github.com/mitchellh/mapstructure"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Filter struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type Global struct {
	ClusterName string `mapstructure:"cluster_name"`
	Domain      string `mapstructure:"domain"`
	Environment string `mapstructure:"environment"`
	Provider    string `mapstructure:"provider"`
	Aws         struct {
		Region            string   `mapstructure:"region"`
		AvailabilityZones []string `mapstructure:"availabilityZones"`
		VpcID             struct {
			Filters Filter `mapstructure:"filters,omitempty"`
			Hook    string `mapstructure:"hook,omitempty"`
		} `mapstructure:"vpc_id"`
		Subnets struct {
			Filters Filter `mapstructure:"filters,omitempty"`
			Hook    string `mapstructure:"hook"`
		} `mapstructure:"subnets"`
		UtilitySubnets struct {
			Filters Filter `mapstructure:"filters,omitempty"`
			Hook    string `mapstructure:"hook"`
		} `mapstructure:"utility-subnets"`
	} `mapstructure:"aws"`
}

type Kubectl struct {
	Name      string `mapstructure:"name"`
	URL       string `mapstructure:"url"`
	Namespace string `mapstructure:"namespace"`
}

type Exec struct {
	Command string `mapstructure:"command"`
	Rundir  string `mapstructure:"rundir,omitempty"`
	Env     []struct {
		Name  string `mapstructure:"name"`
		Value string `mapstructure:"value"`
	} `mapstructure:"env"`
}

type Kops struct {
	Name                 string `mapstructure:"name,omitempty"`
	State                string `mapstructure:"state,omitempty"`
	AdminAccess          string `mapstructure:"admin-access,omitempty"`
	APILoadbalancerType  string `mapstructure:"api-loadbalancer-type,omitempty"`
	APISslCertificate    string `mapstructure:"api-ssl-certificate,omitempty"`
	AssociatePublicIP    string `mapstructure:"associate-public-ip,omitempty"`
	Authorization        string `mapstructure:"authorization,omitempty"`
	Bastion              string `mapstructure:"bastion,omitempty"`
	Channel              string `mapstructure:"channel,omitempty"`
	Cloud                string `mapstructure:"cloud,omitempty"`
	CloudLabels          string `mapstructure:"cloud-labels,omitempty"`
	DNS                  string `mapstructure:"dns,omitempty"`
	DNSZone              string `mapstructure:"dns-zone,omitempty"`
	DryRun               string `mapstructure:"dry-run,omitempty"`
	EncryptEtcdStorage   string `mapstructure:"encrypt-etcd-storage,omitempty"`
	Image                string `mapstructure:"image,omitempty"`
	KubernetesVersion    string `mapstructure:"kubernetes-version,omitempty"`
	MasterCount          int    `mapstructure:"master-count,omitempty"`
	MasterPublicName     string `mapstructure:"master-public-name,omitempty"`
	MasterSecurityGroups string `mapstructure:"master-security-groups,omitempty"`
	MasterSize           string `mapstructure:"master-size,omitempty"`
	MasterTenancy        string `mapstructure:"master-tenancy,omitempty"`
	MasterVolumeSize     string `mapstructure:"master-volume-size,omitempty"`
	MasterZones          string `mapstructure:"master-zones,omitempty"`
	Model                string `mapstructure:"model,omitempty"`
	NetworkCidr          string `mapstructure:"network-cidr,omitempty"`
	Networking           string `mapstructure:"networking,omitempty"`
	NodeCount            int    `mapstructure:"node-count,omitempty"`
	NodeSecurityGroups   string `mapstructure:"node-security-groups,omitempty"`
	NodeSize             string `mapstructure:"node-size,omitempty"`
	NodeTenancy          string `mapstructure:"node-tenancy,omitempty"`
	NodeVolumeSize       string `mapstructure:"node-volume-size,omitempty"`
	Out                  string `mapstructure:"out,omitempty"`
	Outout               string `mapstructure:"outout,,omitempty"`
	Project              string `mapstructure:"project,omitempty"`
	SSHAccess            string `mapstructure:"ssh-access,omitempty"`
	SSHPublicKey         string `mapstructure:"ssh-public-key,omitempty"`
	Subnets              string `mapstructure:"subnets,omitempty"`
	Target               string `mapstructure:"target,omitempty"`
	Topology             string `mapstructure:"topology,omitempty"`
	UtilitySubnets       string `mapstructure:"utility-subnets,omitempty"`
	Vpc                  string `mapstructure:"vpc,omitempty"`
	Zones                string `mapstructure:"zones,omitempty"`
	Snippets             struct {
		Cluster string `mapstructure:"cluster,omitempty"`
		Node    string `mapstructure:"node,omitempty"`
		Master  string `mapstructure:"master,omitempty"`
	} `mapstructure:"snippets,omitempty"`
}

type ConfigData struct {
	Global   Global
	Kubectl  []Kubectl
	Helm     interface{}
	Kops     Kops
	Exec     []Exec
	PreExec  []Exec
	PostExec []Exec
}

/* global holding all the yaml config */
var config ConfigData
var helmfileCmd string = "helmfile"
var helmCmd string = "helm"
var kubectlCmd string = "kubectl"
var kopsCmd string = "kops"
var dockerCmd string = "docker"

func CmdLookPath(cmd string, required ...bool) string {
	var isRequired bool
	if len(required) > 0 {
		isRequired = required[0]
	} else {
		isRequired = false
	}

	path, err := exec.LookPath(cmd)
	if err != nil && isRequired {
		log.Fatalf("Could not find %s in the path\n", cmd)
		return ""
	} else {
		log.Printf("using %s\n", path)
	}
	return path
}

func LoadConfig(configFile string) {
	var dockerRequired = false
	var kubeCmdRequired = true

	if dockerImg != "" {
		dockerRequired = true
		kubeCmdRequired = false
	}
	dockerCmd = CmdLookPath("docker", dockerRequired)
	helmfileCmd = CmdLookPath("helmfile", kubeCmdRequired)
	helmCmd = CmdLookPath("helm", kubeCmdRequired)
	kubectlCmd = CmdLookPath("kubectl", kubeCmdRequired)
	kopsCmd = CmdLookPath("kops", kubeCmdRequired)

	log.Printf("Loading config file from %s\n", configFile)

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error opening %s: #%v\n", configFile, err)
		os.Exit(1)
	}

	c := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	err = mapstructure.Decode(c["kops"], &config.Kops)
	if err != nil {
		log.Fatalf("Error decoding global section: %v\n", err)
	}

	err = mapstructure.Decode(c["global"], &config.Global)
	if err != nil {
		log.Fatalf("Error decoding global section: %v\n", err)
	}

	if HasKey(c, "kubectl") {
		err = mapstructure.Decode(c["kubectl"], &config.Kubectl)
		if err != nil {
			log.Fatalf("Error decoding global section: %v\n", err)
		}
	}

	if HasKey(c, "exec") {
		err = mapstructure.Decode(c["exec"], &config.Exec)
		if err != nil {
			log.Printf("Error decoding exec section: %v; ignoring\n", err)
		}
	}

	if HasKey(c, "pre_exec") {
		err = mapstructure.Decode(c["pre_exec"], &config.PreExec)
		if err != nil {
			log.Printf("Error decoding pre_exec section: %v; ignoring\n", err)
		}
	}

	if HasKey(c, "post_exec") {
		err = mapstructure.Decode(c["post_exec"], &config.PostExec)
		if err != nil {
			log.Printf("Error decoding post_exec section: %v; ignoring\n", err)
		}
	}

	// Helm config
	if HasKey(c, "helm") {
		config.Helm = c["helm"]
	}
}
