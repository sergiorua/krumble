package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const encryptionConfig = `
kind: EncryptionConfig
apiVersion: v1
resources:
  - resources:
    - secrets
    providers:
    - aescbc:
        keys:
        - name: key1
          secret: %s
    - identity: {}
`

// kops create secret encryptionconfig -f $tmpfile --name $CLUSTER_NAME
func genEncryptionConfig() error {
	r := RandStringBytes(32)
	encoded := base64.StdEncoding.EncodeToString([]byte(r))

	tmpfile, err := ioutil.TempFile(tempDir, "kops.*.yaml")
	if err != nil {
		log.Printf("Cannot create temp file for encryption: %v\n", err)
		return err
	}

	c := fmt.Sprintf(encryptionConfig, encoded)
	if debug {
		log.Println(c)
	}
	if _, err := tmpfile.Write([]byte(c)); err != nil {
		log.Printf("Cannot write to temp file: %v\n", err)
		return err
	}
	tmpfile.Close()

	args := []string{"create", "secret", "encryptionconfig", "-f", tmpfile.Name(), "--name", config.Kops.Name, "--state", config.Kops.State}
	err = runCommand(kopsCmd, args...)
	if err != nil {
		return err
	}

	return nil
}

func subnets2string() string {
	var res []string
	for s := range subnets {
		res = append(res, subnets[s].SubnetId)
	}
	return strings.Join(res, ",")
}

func utility_subnets2string() string {
	var res []string
	for s := range utility_subnets {
		res = append(res, utility_subnets[s].SubnetId)
	}
	return strings.Join(res, ",")
}

func BuildKopsCommand() []string {
	cmd := []string{"create", "cluster",
		"--name", config.Kops.Name,
		"--state", config.Kops.State}

	if config.Kops.AdminAccess != "" {
		cmd = append(cmd, "--admin-access", config.Kops.AdminAccess)
	}
	if config.Kops.MasterSize != "" {
		cmd = append(cmd, "--master-size", config.Kops.MasterSize)
	}
	if config.Kops.MasterCount > 0 {
		cmd = append(cmd, fmt.Sprintf("--master-count=%d", config.Kops.MasterCount))
	} else {
		cmd = append(cmd, "--master-count=1")
	}
	if config.Kops.NodeCount > 0 {
		cmd = append(cmd, fmt.Sprintf("--node-count=%d", config.Kops.NodeCount))
	} else {
		cmd = append(cmd, "--node-count", "1")
	}
	if config.Kops.Vpc != "" {
		cmd = append(cmd, "--vpc", string(config.Kops.Vpc))
	} else {
		cmd = append(cmd, "--vpc", vpc.VpcId)
	}
	if config.Kops.Subnets != "" {
		cmd = append(cmd, "--subnets", config.Kops.Subnets)
	} else {
		cmd = append(cmd, "--subnets", subnets2string())
	}
	if config.Kops.KubernetesVersion != "" {
		cmd = append(cmd, "--kubernetes-version", config.Kops.KubernetesVersion)
	}
	if config.Kops.Topology != "" {
		cmd = append(cmd, "--topology", config.Kops.Topology)
	}
	if config.Kops.APISslCertificate != "" {
		cmd = append(cmd, "--api-ssl-certificate", config.Kops.APISslCertificate)
	}
	if config.Kops.Image != "" {
		cmd = append(cmd, "--image", config.Kops.Image)
	}
	if config.Kops.Networking != "" {
		cmd = append(cmd, "--networking", config.Kops.Networking)
	}
	if config.Kops.Cloud != "" {
		cmd = append(cmd, "--cloud", config.Kops.Cloud)
	} else {
		cmd = append(cmd, "--cloud", "aws")
	}
	if config.Kops.SSHPublicKey != "" {
		cmd = append(cmd, "--ssh-public-key", config.Kops.SSHPublicKey)
	}
	if config.Kops.MasterPublicName != "" {
		cmd = append(cmd, "--master-public-name", config.Kops.MasterPublicName)
	} else {
		config.Kops.MasterPublicName = fmt.Sprintf("https://api-%s", config.Kops.Name)
	}
	if config.Kops.MasterZones != "" {
		cmd = append(cmd, "--master-zones", config.Kops.MasterZones)
	} else {
		cmd = append(cmd, "--master-zones", strings.Join(config.Global.Aws.AvailabilityZones, ","))
	}
	if config.Kops.Zones != "" {
		cmd = append(cmd, "--zones", config.Kops.Zones)
	} else {
		cmd = append(cmd, "--zones", strings.Join(config.Global.Aws.AvailabilityZones, ","))
	}
	if config.Kops.UtilitySubnets != "" {
		cmd = append(cmd, "--utility_subnets", config.Kops.UtilitySubnets)
	} else {
		cmd = append(cmd, "--utility-subnets", utility_subnets2string())
	}

	return cmd
}

func RunKops() error {
	os.Setenv("KOPS_STATE_STORE", config.Kops.State)
	os.Setenv("KOPS_CLUSTER_NAME", config.Kops.Name)

	full_cmd := BuildKopsCommand()

	err := runCommand(kopsCmd, full_cmd...)
	if err != nil {
		return err
	}

	genEncryptionConfig()
	mergeKopsConfigs()

	/* now kick off the build */
	// kops update cluster $CLUSTER_NAME --yes
	full_cmd = []string{"update", "cluster", "--yes", "--name", config.Kops.Name, "--state", config.Kops.State}
	err = runCommand(kopsCmd, full_cmd...)
	return err
}

func mergeKopsConfigs() {
	if exists(config.Kops.Snippets.Cluster) {
		log.Printf("Merging from %s\n", config.Kops.Snippets.Cluster)
		MergeKopsClusterSnippets(config.Kops.Snippets.Cluster)
	}
	if exists(config.Kops.Snippets.Node) {
		log.Printf("Merging from %s\n", config.Kops.Snippets.Node)
		MergeKopsNodeSnippets(config.Kops.Snippets.Node)
	}
	if exists(config.Kops.Snippets.Master) {
		log.Printf("Merging from %s\n", config.Kops.Snippets.Master)
		MergeKopsMasterSnippets(config.Kops.Snippets.Master)
	}
}

/*
 * runs kops and then waits until all nodes and masters are running
 *
 */
func ProcessKops() error {
	err := RunKops()
	if err != nil {
		log.Fatal("Error running Kops: %v\n", err)
		return err
	}

	if !dryrun {
		KopsNodesUp()
	}
	return nil
}
