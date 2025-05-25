// Package localaddr provides a simple way to get the local IPv4 address of the machine.
//
// This package is designed as a utility for personal projects where you need to quickly
// retrieve the machine's non-loopback IPv4 address.
package localaddr

import (
	"fmt"
	"net"
)

// Get returns the first non-loopback IPv4 address of an up interface.
//
// It iterates through all network interfaces, skipping those that are down or loopback.
// For each interface, it looks for the first valid IPv4 address and returns it as a string.
//
// Returns:
//   - string: The IPv4 address as a string (e.g., "192.168.1.2")
//   - error: An error if no suitable address is found or if there's an issue accessing network interfaces
func Get() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, v := range interfaces {
		if v.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if v.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := v.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("not connected to the network")
}
