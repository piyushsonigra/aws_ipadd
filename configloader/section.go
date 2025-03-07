package configloader

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

// Constants for AWS IP check URL and configuration file path
const (
	configFilePath = "%s/.aws_ipadd/aws_ipadd"
)

// Get config file from env var or default
func getConfigFile() string {
	// Check if a custom config file path is provided via environment variable
	customPath := os.Getenv("CUSTOM_AWS_IPADD_CONFIG_FILE")
	if customPath != "" {
		// If the custom path is set, return it as-is
		return customPath
	}
	// Fallback to the default configuration path
	return fmt.Sprintf(configFilePath, os.Getenv("HOME"))
}

// Get profile section from config file
func GetSection(profile string) (*ini.Section, error) {
	configFilePath := getConfigFile()
	loadConfig, err := ini.Load(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: \n%v", err)
	}

	// Check if profile exists in config file
	if !loadConfig.HasSection(profile) {
		return nil, fmt.Errorf("profile \"%s\" doesn't exist in \"%s\"", profile, configFilePath)
	}

	section := loadConfig.Section(profile)
	return section, nil
}
