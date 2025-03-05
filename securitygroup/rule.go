package securitygroup

import (
	"aws_ipadd/configloader"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Processes security group rules
func ProcessRule(profile string, ruleConfig *configloader.SecurityGroupRule) (string, error) {
	fmt.Printf("---------------\n%v\n---------------\n", profile)

	// Load AWS SDK configuration with specified profile and region
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(ruleConfig.AWSProfile),
		config.WithRegion(ruleConfig.Region),
	)
	if err != nil {
		return "", err
	}

	// Create EC2 client and retrieve security group details
	ec2Client := ec2.NewFromConfig(awsCfg)
	describeInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{ruleConfig.SecurityGroupID},
	}
	securityGroup, err := ec2Client.DescribeSecurityGroups(context.TODO(), describeInput)
	if err != nil {
		return "", err
	}
	securityRules := securityGroup.SecurityGroups[0].IpPermissions

	// Filter security rules for the specified port and protocol
	var matchingRules []types.IpPermission

	for _, rule := range securityRules {
		protocol := derefString(rule.IpProtocol) // Convert *string to string safely

		switch {
		case protocol == "-1":
			// "all" protocol means all traffic, represented by "-1" in AWS security rules
			if ruleConfig.Protocol == "all" {
				matchingRules = append(matchingRules, rule)
			}

		case protocol == "tcp" || protocol == "udp":
			// Validate TCP and UDP rules against ruleConfig
			if protocol == ruleConfig.Protocol &&
				derefInt(rule.FromPort) == ruleConfig.FromPort &&
				derefInt(rule.ToPort) == ruleConfig.ToPort {

				matchingRules = append(matchingRules, rule)
			}

		default:
			// Handle unsupported protocols
			return "", fmt.Errorf("unable to handle protocol %s, valid values are tcp, udp, all", ruleConfig.Protocol)
		}
	}

	// Define the new security group rule with the current IP
	newRule := types.IpPermission{
		IpProtocol: aws.String(ruleConfig.Protocol),
		FromPort:   aws.Int32(ruleConfig.FromPort),
		ToPort:     aws.Int32(ruleConfig.ToPort),
		IpRanges: []types.IpRange{
			{
				CidrIp:      aws.String(ruleConfig.IP),
				Description: aws.String(ruleConfig.RuleName),
			},
		},
	}

	// Create rule if there is no matching rule for requested port
	if len(matchingRules) == 0 {
		res, err := allowIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, &ruleConfig.IP)
		if err != nil {
			return "", err
		}
		fmt.Println(res)
		return "", nil
	}

	// If the rule with current public IP already exists else create the rule
	for _, matchingRule := range matchingRules {
		if *matchingRule.IpRanges[0].CidrIp == ruleConfig.IP {
			resFmt := fmt.Sprintf("Your IP %s is already whitelisted for FromPort %d to ToPort %d.\n", ruleConfig.IP, ruleConfig.FromPort, ruleConfig.ToPort)
			if ruleConfig.Protocol == "all" {
				resFmt = fmt.Sprintf("Your IP %s is already whitelisted for all traffic.", ruleConfig.IP)
			}
			fmt.Print(resFmt)
			return "", nil
		} else {
			res, err := allowIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, &ruleConfig.IP)
			if err != nil {
				return "", err
			}
			fmt.Println(res)
			return "", nil
		}
	}

	// Remove and Add rule if port and rule name matches with requested rule but whitelisted IP in rule is different
	for _, matchingRule := range matchingRules {
		if derefString(matchingRule.IpRanges[0].Description) == ruleConfig.RuleName && derefString(matchingRule.IpRanges[0].CidrIp) != ruleConfig.IP {
			fmt.Println("Whitelisting your current IP...")
			// Revoke old IP permission
			res, err := revokeIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, *matchingRule.IpRanges[0].CidrIp)
			if err != nil {
				return "", err
			}
			fmt.Println(res)

			// Allow new IP permission
			res, err = allowIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, &ruleConfig.IP)
			if err != nil {
				return "", err
			}
			fmt.Println(res)
		} else {
			return "Nothing to modify", nil
		}
	}
	return "Nothing to process", nil
}

// Helper function to safely dereference an *int pointer
// If the pointer is nil, return 0 to avoid runtime panics
func derefInt(p *int32) int32 {
	if p != nil {
		return *p
	}
	return 0
}

// Helper function to safely dereference a *string pointer
// If the pointer is nil, return an empty string
func derefString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}
