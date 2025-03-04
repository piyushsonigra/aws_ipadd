package cliargs

import (
	"flag"
	"fmt"
	"os"
)

// Args stores command-line arguments.
type Args struct {
	Profile  string
	Port     string
	FromPort string
	ToPort   string
	Protocol string
	IP       string
	RuleName string
}

// Custom usage function for better formatting
func customUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  aws_ipadd --profile <profile-name>")
	fmt.Println("  aws_ipadd --profile <profile-name> --port <port> --current_ip <current_ip> [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  --profile <string>     aws_ipadd profile name (required)")
	fmt.Println("  --port <int>           Port number, this will be ignored if from_port and to_port is passed (optional)")
	fmt.Println("  --from_port <int>      From port number, It should be passed with to_port. Only from_port is not valid argument (optional)")
	fmt.Println("  --to_port <int>        To port number, It should be passed with from_port. Only to_port is not valid argument (optional)")
	fmt.Println("  --protocol <int>       Protocol e.g TCP, UPD, all (optional)")
	fmt.Println("  --ip <string>          IP address with subnetmask e.g '10.10.19.1/32' (optional)")
	fmt.Println("  --rule_name <string>   Security group rule name (optional)")
}

func ParseArgs() *Args {
	args := &Args{}
	flag.Usage = customUsage

	flag.StringVar(&args.Profile, "profile", "", "AWS profile name (required)")
	flag.StringVar(&args.Port, "port", "", "Port number, this will be ignored if from_port and to_port is passed (optional)")
	flag.StringVar(&args.FromPort, "from_port", "", "Port number, It should be passed with to_port. Only from_port is not valid argument (optional)")
	flag.StringVar(&args.ToPort, "to_port", "", "Port number, It should be passed with from_port. Only to_port is not valid argument (optional)")
	flag.StringVar(&args.Protocol, "protocol", "", "Protocol e.g., TCP, UDP, all (optional)")
	flag.StringVar(&args.IP, "ip", "", "IP address with subnetmask e.g., '10.10.19.1/32' (optional)")
	flag.StringVar(&args.RuleName, "rule_name", "", "Security group rule name (optional)")

	flag.Parse()

	if args.Profile == "" {
		fmt.Println("Error: --profile is required")
		flag.Usage()
		os.Exit(1)
	}

	return args
}
