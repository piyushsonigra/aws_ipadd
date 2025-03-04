package main

import (
	"aws_ipadd/cliargs"
	"aws_ipadd/configloader"
	"aws_ipadd/securitygroup"
	"log"
)

func main() {

	// Get CLI args
	args := cliargs.ParseArgs()

	// Load profile config from file
	section, err := configloader.GetSection(args.Profile)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare security group rule
	rule, err := configloader.GetConfig(section, args)
	if err != nil {
		log.Fatal(err)
	}

	// Process security group rule
	_, err = securitygroup.ProcessRule(args.Profile, rule)
	if err != nil {
		log.Fatal(err)
	}
}
