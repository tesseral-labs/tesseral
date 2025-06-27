package restrictedhttp

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestrictedDial(t *testing.T) {
	t.Parallel()

	pass := [...]string{
		"accounts.google.com:443",
		"example.com:443",
	}

	fail := [...]string{
		"0.0.0.0:443",
		"10.10.10.10:443",
		"100.64.0.1:443",
		"127.0.0.1:443",
		"169.254.169.255:443",
		"172.16.42.1:443",
		"192.168.0.1:443",
		"[::1]:443",
		"[fc00::1]:443",
		"[fe80::1]:443",
		"169.254.169.254:80",
	}

	for _, addr := range pass {
		t.Run("trying to connect to "+addr+" must pass", func(t *testing.T) {
			t.Parallel()

			dial := restrictedDial(
				(&net.Dialer{}).DialContext,
			)
			_, err := dial(t.Context(), "tcp", addr)
			require.NoError(t, err)
		})
	}

	for _, addr := range fail {
		t.Run("trying to connect to "+addr+" must fail", func(t *testing.T) {
			t.Parallel()

			dial := restrictedDial(
				(&net.Dialer{}).DialContext,
			)
			_, err := dial(t.Context(), "tcp", addr)
			require.Error(t, err)
		})
	}
}

func TestRestrictedHttp(t *testing.T) {
	t.Parallel()

	pass := [...]string{
		"accounts.google.com",
		"example.com",
	}

	fail := [...]string{
		"0.0.0.0",
		"10.10.10.10",
		"100.64.0.1",
		"127.0.0.1",
		"169.254.169.255",
		"172.16.42.1",
		"192.168.0.1",
		"[::1]:443",
		"[fc00::1]",
		"[fe80::1]",
		"169.254.169.254",
	}

	for _, addr := range pass {
		t.Run("trying to connect to "+addr+" must pass", func(t *testing.T) {
			t.Parallel()

			client := NewClient(nil)

			_, err := client.Get("http://" + addr)
			require.NoError(t, err)

			_, err = client.Get("https://" + addr)
			require.NoError(t, err)
		})
	}

	for _, addr := range fail {
		t.Run("trying to connect to "+addr+" must fail", func(t *testing.T) {
			t.Parallel()

			client := NewClient(nil)

			_, err := client.Get("http://" + addr)
			require.Error(t, err)

			_, err = client.Get("https://" + addr)
			require.Error(t, err)
		})
	}

	t.Run("restricted HTTP client cannot connect to local address", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer srv.Close()

		client := NewClient(nil)
		_, err := client.Get(srv.URL)
		require.Error(t, err)
	})
}
