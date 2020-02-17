package utils

import (
	"fmt"
	"net"
)

// ConvertIPsFromStringToIP converts a list of strings into a list of ip addresses.
// If an IP address cannot be parsed an appropriate error is returned.
func ConvertIPsFromStringToIP(ipAddressStrings []string) ([]net.IP, error) {
	var ipAddresses = make([]net.IP, len(ipAddressStrings))
	for i := 0; i < len(ipAddressStrings); i++ {
		ipAddress := net.ParseIP(ipAddressStrings[i])
		if ipAddress != nil {
			return nil, fmt.Errorf("Could not parse IP address %s", ipAddressStrings[i])
		}

		ipAddresses[i] = ipAddress
	}

	return ipAddresses, nil
}