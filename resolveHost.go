package socks4

import (
	"fmt"
	"net"
)

func resolveHost(host string) (net.IP, error) {

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	// find and return the first valid IPv4
	for _, ip := range ips {
		i := ip.To4()
		if i != nil && len(i) == net.IPv4len {
			return i, nil
		}
	}

	return nil, fmt.Errorf("couldn't resolve hostname: %s", host)

}
