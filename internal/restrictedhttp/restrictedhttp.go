package restrictedhttp

// Code modified from https://github.com/segmentio/netsec

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func NewTransport() *http.Transport {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	dialer := &restrictedDialer{
		dial:     defaultTransport.DialContext,
		denyList: privateIPNetworks,
	}

	// Construct a new transport with the same parameters as http.DefaultTransport
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				slog.ErrorContext(ctx, "restricted_dial_error", "network", network, "addr", addr, "error", err)
			}
			return conn, err
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

type restrictedDialer struct {
	dial     func(ctx context.Context, network, addr string) (net.Conn, error)
	denyList denyList
}

func (d *restrictedDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to split host and port from '%s': %w", addr, err)
	}

	var ipAddr *net.IPAddr
	if ip := net.ParseIP(host); ip != nil {
		ipAddr = &net.IPAddr{IP: ip}
	} else {
		addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve IP for host '%s': %w", host, err)
		}
		if len(addrs) == 0 {
			// I'm not sure this could ever happen but the net package
			// documentation does not ensure that the list of addresses will
			// never be empty if the the error is nil.
			return nil, &net.OpError{
				Op:  "lookup",
				Net: network,
				Err: &net.DNSError{
					Err:  "No addresses returned by the DNS resolver",
					Name: host,
				},
			}
		}
		ipAddr = &addrs[0]
	}

	if err := d.denyList.check(ipAddr); err != nil {
		return nil, &net.OpError{
			Op:   "dial",
			Net:  network,
			Addr: ipAddr,
			Err: &net.AddrError{
				Err:  err.Error(),
				Addr: addr,
			},
		}
	}

	return d.dial(ctx, network, addr)
}
