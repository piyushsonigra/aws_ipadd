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
	fmt.Printf("---------------\n%s\n---------------\nSecurityGroupID: %s\n", profile, ruleConfig.SecurityGroupID)

	// Get security group rules
	ec2Client, securityGroupRules, err := getSecurityGroupRules(ruleConfig)
	if err != nil {
		return "", err
	}

	// Get valid allowed rules
	validRules, err := getValidRules(ruleConfig, &securityGroupRules)
	if err != nil {
		return "", err
	}

	// findMatchingRule searches for matching rules by IP or rule name
	ruleIPMatched, ruleNameMatched, matchingRule := getMatchingRule(ruleConfig, &validRules)

	// Return if IP matched in rule
	if ruleIPMatched {
		resFmt := fmt.Sprintf("Your IP %s is already whitelisted for FromPort %d to ToPort %d.\n", ruleConfig.IP, ruleConfig.FromPort, ruleConfig.ToPort)
		if ruleConfig.Protocol == "all" {
			resFmt = fmt.Sprintf("Your IP %s is already whitelisted for all traffic.\n", ruleConfig.IP)
		}
		fmt.Print(resFmt)
		return "", nil
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
	if !ruleIPMatched && !ruleNameMatched {
		res, err := allowIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, &ruleConfig.IP)
		if err != nil {
			return "", err
		}
		fmt.Println(res)
		return "", nil
	}

	// Modify security group rule
	if ruleNameMatched {
		// Revoke old IP permission
		if ruleConfig.RuleName != "" {
			fmt.Println("Updating your current IP...")
			res, err := revokeIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, *matchingRule.IpRanges[0].CidrIp)
			if err != nil {
				return "", err
			}
			fmt.Println(res)
		}

		// Allow new IP permission
		res, err := allowIPPermission(ec2Client, &ruleConfig.SecurityGroupID, newRule, &ruleConfig.IP)
		if err != nil {
			return "", err
		}
		fmt.Println(res)

		return "", nil
	}

	fmt.Println("Nothing to process")
	return "", nil
}

// Get security group rules list
func getSecurityGroupRules(ruleConfig *configloader.SecurityGroupRule) (*ec2.Client, []types.IpPermission, error) {

	// Load AWS SDK configuration with specified profile and region
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(ruleConfig.AWSProfile),
		config.WithRegion(ruleConfig.Region),
	)
	if err != nil {
		return nil, nil, err
	}

	// Create EC2 client and retrieve security group details
	ec2Client := ec2.NewFromConfig(awsCfg)
	describeInput := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{ruleConfig.SecurityGroupID},
	}
	securityGroup, err := ec2Client.DescribeSecurityGroups(context.TODO(), describeInput)
	if err != nil {
		return nil, nil, err
	}
	return ec2Client, securityGroup.SecurityGroups[0].IpPermissions, nil
}

// Filter Valid rules for the allowed port and protocol
func getValidRules(ruleConfig *configloader.SecurityGroupRule, securityGroupRules *[]types.IpPermission) ([]types.IpPermission, error) {
	var matchingRules []types.IpPermission
	for _, rule := range *securityGroupRules {
		protocol := derefString(rule.IpProtocol) // Convert *string to string safely

		switch {
		case protocol == "-1":
			// "all" protocol means all traffic, represented by "-1" in AWS security rules
			if ruleConfig.Protocol == "all" {
				matchingRules = append(matchingRules, rule)
				return matchingRules, nil
			}

		case protocol == "tcp" || protocol == "udp":
			// Validate TCP and UDP rules against ruleConfig
			if protocol == ruleConfig.Protocol &&
				derefInt(rule.FromPort) == ruleConfig.FromPort &&
				derefInt(rule.ToPort) == ruleConfig.ToPort {

				matchingRules = append(matchingRules, rule)
				return matchingRules, nil
			}

		default:
			return nil, fmt.Errorf("unable to handle protocol %s, valid values are tcp, udp, all", ruleConfig.Protocol)
		}
	}
	return matchingRules, nil
}

// Get matching rule for IP and Rule name
func getMatchingRule(ruleConfig *configloader.SecurityGroupRule, validRules *[]types.IpPermission) (bool, bool, types.IpPermission) {
	var matchingRule types.IpPermission
	var ruleIPMatched, ruleNameMatched bool

	for _, rule := range *validRules {
		for _, ipRange := range rule.IpRanges {
			matchingRule = rule
			if derefString(ipRange.CidrIp) == ruleConfig.IP {
				ruleIPMatched = true
				matchingRule.IpRanges = []types.IpRange{ipRange}
			}
			if derefString(ipRange.Description) == ruleConfig.RuleName {
				ruleNameMatched = true
				matchingRule.IpRanges = []types.IpRange{ipRange}
			}
		}
		if ruleIPMatched || ruleNameMatched {
			break
		}
	}
	return ruleIPMatched, ruleNameMatched, matchingRule
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
