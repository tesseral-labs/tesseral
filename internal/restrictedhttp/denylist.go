package restrictedhttp

import (
	"fmt"
	"net"
)

type denyList []*net.IPNet

func (d denyList) check(addr net.Addr) error {
	ip := ipAddressOf(addr)
	for _, ipnet := range d {
		if ipnet.Contains(ip) {
			return fmt.Errorf("unauthorized attempt to connect to an address in a denied network (%s)", addr)
		}
	}
	return nil
}

var (
	// privateIPNetworks lists standard IP networks used for private networks.
	privateIPNetworks = denyList{
		cidr("0.0.0.0/32"),
		cidr("10.0.0.0/8"),
		cidr("100.64.0.0/10"),
		cidr("127.0.0.0/8"),
		cidr("169.254.0.0/16"),
		cidr("172.16.0.0/12"),
		cidr("192.168.0.0/16"),
		cidr("fc00::/7"),
		cidr("fd00::/8"),
		cidr("fe80::/10"),
		cidr("::1/128"),
	}
)

// cidr is like net.ParseCIDR but panics if the input is invalid. This function
// is useful to initialize lists of CIDRs without having to check errors.
func cidr(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return ipnet
}

// ipAddressOf extracts the IP address of addr, or returns nil if none were
// found.
func ipAddressOf(addr net.Addr) net.IP {
	switch a := addr.(type) {
	case *net.TCPAddr:
		return a.IP
	case *net.UDPAddr:
		return a.IP
	case *net.IPAddr:
		return a.IP
	default:
		return nil
	}
}
