package securitygroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// RevokeIPPermission removes an existing security group rule
func revokeIPPermission(client *ec2.Client, securityGroupID *string, rule types.IpPermission, oldIP string) (string, error) {

	// Update rule with the new IP to be allowed
	rule.IpRanges[0].CidrIp = &oldIP
	input := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       securityGroupID,
		IpPermissions: []types.IpPermission{rule},
	}

	res, err := client.RevokeSecurityGroupIngress(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to revoke IP: %v", err)
	}

	// Handling UnknownIpPermissions
	if len(res.UnknownIpPermissions) > 0 {
		fmt.Println("Unknown IP Permissions detected. The following rules were not found:")
		for _, perm := range res.UnknownIpPermissions {
			fmt.Printf("     Protocol: %s, Port: %d-%d", aws.ToString(perm.IpProtocol), aws.ToInt32(perm.FromPort), aws.ToInt32(perm.ToPort))
			for _, ipRange := range perm.IpRanges {
				fmt.Printf(", CIDR: %s, Description: %s\n", aws.ToString(ipRange.CidrIp), aws.ToString(ipRange.Description))
			}
		}
		return "", fmt.Errorf("failed")
	}

	resFmt := fmt.Sprintf("Removed old whitelisted IP %s for FromPort %d to ToPort %d", oldIP, *rule.FromPort, *rule.ToPort)
	if *rule.FromPort == 0 && *rule.ToPort == 0 {
		resFmt = fmt.Sprintf("Removed old whitelisted IP %s for all traffic", oldIP)
	}

	return resFmt, nil
}
