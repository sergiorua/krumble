package cmd

import (
	"fmt"
	"log"
	"strings"
)

func subnets2string() string {
	var res []string
	for s := range subnets {
		res = append(res, subnets[s].SubnetId)
	}
	return strings.Join(res, ",")
}

func BuildKopsCommand() []string {
	cmd := []string{"--name", config.Kops.Name,
		"--state", config.Kops.State}

	if config.Kops.AdminAccess != "" {
		cmd = append(cmd, "--admin-access", config.Kops.AdminAccess)
	}
	if config.Kops.MasterSize != "" {
		cmd = append(cmd, "--master-size", config.Kops.MasterSize)
	}
	if config.Kops.MasterCount > 0 {
		cmd = append(cmd, fmt.Sprintf("--master-count %d", config.Kops.MasterCount))
	} else {
		cmd = append(cmd, "--master-count", "1")
	}
	if config.Kops.NodeCount > 0 {
		cmd = append(cmd, fmt.Sprintf("--node-count %d", config.Kops.NodeCount))
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
	if config.Kops.UtilitySubnets != "" {
		cmd = append(cmd, "--utility-subnets", config.Kops.UtilitySubnets)
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
	if config.Kops.MasterZones != "" {
		cmd = append(cmd, "--master-zones", config.Kops.MasterZones)
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

	return cmd
}

func ProcessKops() error {
	full_cmd := BuildKopsCommand()

	if dryrun {
		log.Printf("CMD: %s %s\n", kopsCmd, full_cmd)
		return nil
	}

	err := runCommand(kopsCmd, full_cmd...)
	return err
}
