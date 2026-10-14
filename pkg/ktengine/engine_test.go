package ktengine

import (
	"github.com/cli/cli/v2/internal/config"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestCredentialFreeConfigIgnoresOperatorEnvironment(t *testing.T) {
	t.Setenv("GH_TOKEN", "must-not-be-loaded")
	t.Setenv("GH_CONFIG_DIR", t.TempDir())
	cfg := config.NewTransportConfig()
	token, source := cfg.Authentication().ActiveToken("github.com")
	require.Empty(t, token)
	require.Equal(t, "transport", source)
	require.False(t, cfg.Authentication().HasEnvToken())
	require.Error(t, cfg.Write())
}

func TestOfficialCommandUsesInjectedHTTPWithoutToken(t *testing.T) {
	t.Setenv("GH_TOKEN", "must-not-be-loaded")
	ios, _, out, _ := iostreams.Test()
	called := false
	c := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		require.Equal(t, "api.github.com", r.URL.Host)
		require.Equal(t, "/meta", r.URL.Path)
		require.Empty(t, r.Header.Get("Authorization"))
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"verifiable_password_authentication":true}`)), Request: r}, nil
	})}
	cmd, e := New(ios, c, c, "")
	require.NoError(t, e)
	cmd.SetArgs([]string{"api", "meta"})
	require.NoError(t, cmd.Execute())
	require.True(t, called)
	require.Contains(t, out.String(), "verifiable_password_authentication")
}
