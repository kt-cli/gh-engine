package ktcodespaces

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestOpenAuthorizesActualResourceBeforeTunnelRequests(t *testing.T) {
	called := 0
	client := &http.Client{Transport: testTransport(func(r *http.Request) (*http.Response, error) {
		called++
		require.Equal(t, "true", r.URL.Query().Get("internal"))
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{"name":"fixture","state":"Available","owner":{"login":"account"},"repository":{"full_name":"acme/widgets","owner":{"login":"acme"}},"connection":{"tunnelProperties":{"serviceUri":"https://untrusted.invalid"}}}`))}, nil
	})}
	denied := errors.New("fixture denial")
	_, e := Open(context.Background(), "fixture", client, func(cs Codespace) error {
		require.Equal(t, "account", cs.Account)
		require.Equal(t, "acme", cs.Owner)
		require.Equal(t, "widgets", cs.Repository)
		return denied
	}, func(string, string) error {
		t.Fatal("tunnel destination used before resource authorization")
		return nil
	})
	require.ErrorIs(t, e, denied)
	require.Equal(t, 1, called)
}
