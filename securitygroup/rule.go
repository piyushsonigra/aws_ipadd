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

	// Filter security rules for the specified port
	var matchingRules []types.IpPermission
	for _, rule := range securityRules {
		if ruleConfig.Protocol == "all" && *rule.IpProtocol == "-1" {
			matchingRules = append(matchingRules, rule)
		}
		if ruleConfig.Protocol != "all" &&
			*rule.FromPort == ruleConfig.FromPort &&
			*rule.ToPort == ruleConfig.ToPort &&
			*rule.IpProtocol == ruleConfig.Protocol {

			matchingRules = append(matchingRules, rule)
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

	// Complete the script if the rule with current public IP already exists
	for _, matchingRule := range matchingRules {
		if *matchingRule.IpRanges[0].CidrIp == ruleConfig.IP {
			resFmt := fmt.Sprintf("Your IP %s is already whitelisted for FromPort %d to ToPort %d.\n", ruleConfig.IP, ruleConfig.FromPort, ruleConfig.ToPort)
			if ruleConfig.Protocol == "all" {
				resFmt = fmt.Sprintf("Your IP %s is already whitelisted for all traffic.", ruleConfig.IP)
			}
			fmt.Println(resFmt)
			return "", nil
		}
	}

	// Remove and Add rule if port and rule name matches with requested rule but whitelisted IP in rule is different
	for _, matchingRule := range matchingRules {
		if *matchingRule.IpRanges[0].Description == ruleConfig.RuleName && *matchingRule.IpRanges[0].CidrIp != ruleConfig.IP {
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
		}
	}
	return "", nil
}
