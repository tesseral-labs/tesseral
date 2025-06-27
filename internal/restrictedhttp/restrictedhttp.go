package restrictedhttp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

func NewClient(httpClient *http.Client) *http.Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		transport = http.DefaultTransport.(*http.Transport)
	}
	transport.DialContext = restrictedDial(transport.DialContext)

	return &http.Client{
		Transport:     transport,
		Jar:           httpClient.Jar,
		Timeout:       httpClient.Timeout,
		CheckRedirect: httpClient.CheckRedirect,
	}
}

type dialFunc = func(ctx context.Context, network string, addr string) (net.Conn, error)

func restrictedDial(dial dialFunc) dialFunc {
	dialer := &restrictedDialer{
		dial:     dial,
		denyList: privateIPNetworks,
	}
	return dialFunc(func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			slog.ErrorContext(ctx, "restricted_dial_error", "network", network, "addr", addr, "error", err)
		}
		return conn, err
	})
}

type restrictedDialer struct {
	dial     dialFunc
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
