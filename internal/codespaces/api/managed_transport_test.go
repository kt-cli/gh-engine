package api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cli/cli/v2/internal/codespaces/transport"
	"github.com/stretchr/testify/require"
)

type managedTestTransport struct{ request *http.Request }

func (m *managedTestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	m.request = r
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"name":"fixture","state":"Available"}`)), Header: http.Header{}}, nil
}
func (m *managedTestTransport) ConnectCodespace(context.Context, string) (transport.PortForwarder, error) {
	return nil, nil
}

func TestManagedTransportGetsMetadataWithoutConnectionCredentials(t *testing.T) {
	m := &managedTestTransport{}
	client := &http.Client{Transport: m}
	a := &API{client: func() (*http.Client, error) { return client, nil }, externalClient: func() (*http.Client, error) { return client, nil }, githubAPI: defaultAPIURL}
	cs, e := a.GetCodespace(context.Background(), "fixture", true)
	require.NoError(t, e)
	require.Equal(t, "fixture", cs.Name)
	require.Empty(t, m.request.URL.RawQuery)
}
