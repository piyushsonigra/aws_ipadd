package securitygroup

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Allows a new IP permission in the security group
func allowIPPermission(client *ec2.Client, securityGroupID *string, rule types.IpPermission, newIP *string) (string, error) {

	// Update rule with the new IP to be allowed
	rule.IpRanges[0].CidrIp = newIP
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       securityGroupID,
		IpPermissions: []types.IpPermission{rule},
	}

	// Execute the allow operation
	_, err := client.AuthorizeSecurityGroupIngress(context.TODO(), input)
	if err != nil {
		return "", err
	}

	resFmt := fmt.Sprintf("Whitelisted your IP %s for FromPort %d to ToPort %d.", *rule.IpRanges[0].CidrIp, *rule.FromPort, *rule.ToPort)
	if *rule.FromPort == 0 && *rule.ToPort == 0 {
		resFmt = fmt.Sprintf("Whitelisted your IP %s for all traffic.", *rule.IpRanges[0].CidrIp)
	}

	return resFmt, nil
}
