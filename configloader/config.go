package configloader

import (
	"aws_ipadd/cliargs"
	"aws_ipadd/publicip"
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/ini.v1"
)

// SecurityGroupRule represents the security group rule details extracted from the config file
type SecurityGroupRule struct {
	AWSProfile      string
	Region          string
	SecurityGroupID string
	Protocol        string
	FromPort        int32
	ToPort          int32
	IP              string
	RuleName        string
}

// GetSecurityGroupRule extracts security group rule details from the config file section
func GetConfig(section *ini.Section, args *cliargs.Args) (*SecurityGroupRule, error) {

	rule := &SecurityGroupRule{}

	requiredKeys := []string{"aws_profile", "region_name", "security_group_id"}
	// Validate required keys
	for _, key := range requiredKeys {
		value := section.Key(key).String()
		if value == "" {
			return nil, fmt.Errorf("%s is missing in config file", key)
		}

		// Assign values dynamically based on key
		switch key {
		case "aws_profile":
			rule.AWSProfile = value
		case "region":
			rule.Region = value
		case "security_group_id":
			rule.SecurityGroupID = value
		}
	}

	var err error

	// Get Protocol
	rule.Protocol, err = getProtocolValue(section, args)
	if err != nil {
		return nil, err
	}

	// Get Port value
	rule.FromPort, rule.ToPort, err = getPortValue(section, args, rule.Protocol)
	if err != nil {
		return nil, err
	}

	// Get IP value
	rule.IP, err = getIPvalue(args)
	if err != nil {
		return nil, err
	}

	// Set rule name
	if args.RuleName != "" {
		rule.RuleName = args.RuleName
	} else {
		rule.RuleName = section.Key("rule_name").String()
	}

	return rule, nil
}

// port value
func getPortValue(section *ini.Section, args *cliargs.Args, protocol string) (int32, int32, error) {
	// Extract values from config and CLI arguments
	cfgPort := section.Key("port").String()
	cfgFromPort := section.Key("from_port").String()
	cfgToPort := section.Key("to_port").String()
	argPort := args.Port
	argFromPort := args.FromPort
	argToPort := args.ToPort

	// Check for all traffic rule
	// For Protocol value "all", port is not required to pass, all protocols and ports are allowed for whitelisted IP
	if protocol == "all" {
		return 0, 0, nil
	}

	// Ensure at least one of the port values is provided
	if cfgPort == "" && cfgFromPort == "" && cfgToPort == "" && argPort == "" && argFromPort == "" && argToPort == "" {
		return 0, 0, errors.New("port or from_port and to_port must be provided either in the config file or CLI arguments")
	}

	// Case 1: Use from_port and to_port from CLI arguments if both are provided
	if argFromPort != "" && argToPort != "" {
		fromPort, err := parsePort(argFromPort)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid 'from_port' value in CLI arguments: %v", err)
		}
		toPort, err := parsePort(argToPort)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid 'to_port' value in CLI arguments: %v", err)
		}
		return fromPort, toPort, nil
	}

	// Case 2: Use port from CLI arguments if provided
	if argPort != "" {
		port, err := parsePort(argPort)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid 'port' value in CLI arguments: %v", err)
		}
		return port, port, nil
	}

	// Case 3: Use from_port and to_port from config file if both exist
	if cfgFromPort != "" && cfgToPort != "" {
		fromPort, err := parsePort(cfgFromPort)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid 'from_port' value in config file: %v", err)
		}
		toPort, err := parsePort(cfgToPort)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid 'to_port' value in config file: %v", err)
		}
		return fromPort, toPort, nil
	}

	// Case 4: Use port from config file if available
	if cfgPort != "" {
		port, err := parsePort(cfgPort)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid 'port' value in config file: %v", err)
		}
		return port, port, nil
	}

	// This should never be reached, but return an error just in case
	return 0, 0, errors.New("unexpected error determining port values")
}

// parsePort converts a port string to an int32 value
func parsePort(portStr string) (int32, error) {
	port, err := strconv.ParseInt(portStr, 10, 32)
	return int32(port), err
}

// IP value
func getIPvalue(args *cliargs.Args) (string, error) {

	// Check if IP exist in cli arguments
	ip := args.IP
	var err error
	if ip == "" {
		ip, err = publicip.GetCurrentPublicIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return ip, nil
}

func getProtocolValue(section *ini.Section, args *cliargs.Args) (string, error) {
	if section.Key("protocol").String() == "" && args.Protocol == "" {
		return "", errors.New("protocol value is missing in configfile and cli arguments")
	}
	protocol := section.Key("protocol").String()
	if args.Protocol != "" {
		return args.Protocol, nil
	}
	return protocol, nil
}
