package publicip

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	awsIPCheckURL = "https://checkip.amazonaws.com"
)

// Retrieves the current public IP address of the machine
func GetCurrentPublicIP() (string, error) {
	// Query AWS IP check service
	resp, err := http.Get(awsIPCheckURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and format response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/32", strings.TrimSpace(string(body))), nil
}
