package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os/exec"
)

type AwsVpc struct {
	VpcId     string
	CidrBlock string
}

type AwsSubnet struct {
	SubnetId         string
	CidrBlock        string
	VpcId            string
	AvailabilityZone string
}

func AwsSubnetsLookup(vpcId string, filter Filter) []AwsSubnet {
	var subnets []AwsSubnet

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		log.Fatal("Cannot connec to AWS: %v\n", err)
		return subnets
	}

	svc := ec2.New(sess)
	params := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(filter.Key),
				Values: []*string{aws.String(filter.Value)},
			}, {
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	}
	resp, err := svc.DescribeSubnets(params)
	if err != nil {
		log.Printf("Error discovering Subnets: %v\n", err)
		return subnets
	}

	for _, sb := range resp.Subnets {
		var s AwsSubnet
		s.VpcId = *sb.VpcId
		s.SubnetId = *sb.SubnetId
		s.AvailabilityZone = *sb.AvailabilityZone
		s.CidrBlock = *sb.CidrBlock
		subnets = append(subnets, s)
	}
	return subnets
}

func AwsVpcLookup(filter Filter) AwsVpc {
	var vpc AwsVpc

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		log.Fatal("Cannot connec to AWS: %v\n", err)
		return vpc
	}

	svc := ec2.New(sess)

	params := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String(filter.Key),
				Values: []*string{aws.String(filter.Value)},
			},
		},
	}
	resp, err := svc.DescribeVpcs(params)
	if err != nil {
		log.Printf("Error discovering VPC: %v\n", err)
		return vpc
	}

	vpc.VpcId = *resp.Vpcs[0].VpcId
	vpc.CidrBlock = *resp.Vpcs[0].CidrBlock

	return vpc
}

func RunHook(hook string) string {
	out, err := exec.Command(hook).Output()
	if err != nil {
		log.Fatal("Cmd %v failed: %v\n", hook, err)
		return ""
	}
	return string(out)
}

func ProcessHooks() (AwsVpc, []AwsSubnet) {
	vpc := AwsVpcLookup(config.Global.Aws.VpcID.Filters)
	subnets := AwsSubnetsLookup(vpc.VpcId, config.Global.Aws.Subnets.Filters)

	return vpc, subnets
}