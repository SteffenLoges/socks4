package socks4

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	SOCKS4 = iota
	SOCKS4A
)

func Dialer(protocol int, proxyAddr, userID string, timeout time.Duration) func(string, string) (net.Conn, error) {
	return func(_, targetAddr string) (net.Conn, error) {
		return dialer(protocol, proxyAddr, targetAddr, userID, timeout)
	}
}

func FasthttpDialer(protocol int, proxyAddr, userID string, timeout time.Duration) func(string) (net.Conn, error) {
	return func(targetAddr string) (net.Conn, error) {
		return dialer(protocol, proxyAddr, targetAddr, userID, timeout)
	}
}

func dialer(protocol int, proxyAddr, targetAddr, userID string, timeout time.Duration) (net.Conn, error) {

	// Inintial non-zero target-ip for SOCKS4A proxies
	// Will get replaced later if protocol is SOCKS4
	targetIP := net.IPv4(0, 0, 0, 1)

	targetHost, portStr, err := net.SplitHostPort(targetAddr)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	// Connect to the proxy
	conn, err := net.DialTimeout("tcp", proxyAddr, timeout)
	if err != nil {
		return nil, err
	}

	// Check if destination is IPv6
	if strings.Contains(targetHost, ":") {
		return nil, errors.New("The socks4 protocol doesn't support IPv6")
	}

	// targetHostIP := net.ParseIP(targetHost)
	if protocol == SOCKS4 {
		targetIP, err = resolveHost(targetHost)
		if err != nil {
			return nil, err
		}
	}

	// Building the request
	req := []byte{}

	// SOCKS version number, 0x04 for socks4
	req = append(req, 0x04)

	// command code, 0x01 = establish a TCP/IP stream connection
	req = append(req, 0x01)

	// 2-byte big-endian port number
	req = append(req, byte(port>>8))
	req = append(req, byte(port))

	// 4 byte big-endian ip address
	req = append(req, targetIP.To4()...)

	// the user ID string, null-terminated.
	req = append(req, []byte(userID)...)
	req = append(req, 0x00)

	// append the domain name of the host if we're using SOCKS4A
	// and have a valid target domain, null-terminated
	if protocol == SOCKS4A {
		req = append(req, []byte(targetHost)...)
		req = append(req, 0x00)
	}

	// Send request
	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}

	// 8 byte response
	resp := make([]byte, 8)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, err
	}

	if n != 8 {
		return nil, errors.New("Proxy returned unexpected response packet")
	}

	// resp[1] = reply code
	switch resp[1] {

	case 0x5A:
		break // Request granted

	case 0x5B:
		return nil, errors.New("request rejected or failed")

	case 0x5C:
		return nil, errors.New("request rejected becasue SOCKS server cannot connect to identd on the client")

	case 0x5D:
		return nil, errors.New("request rejected because the client program and identd report different user-ids")

	default:
		return nil, errors.New("Request failed because of an unknown error")

	}

	return conn, nil

}
