// Package ktcodespaces exposes the command core's managed connection seam.
package ktcodespaces

import (
	"context"
	"fmt"
	"github.com/cli/cli/v2/internal/codespaces/api"
	"github.com/cli/cli/v2/internal/codespaces/connection"
	"github.com/cli/cli/v2/internal/codespaces/portforwarder"
	"github.com/cli/cli/v2/internal/codespaces/transport"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"net/http"
	"strings"
)

// PortForwarder is the command core's port forwarding interface.
type PortForwarder = transport.PortForwarder

// ForwardPortOpts selects a remote port and forwarding behavior.
type ForwardPortOpts = transport.ForwardPortOpts

// Codespace identifies the actual upstream resource without connection credentials.
type Codespace struct{ Name, Account, Owner, Repository, State string }

// Open connects using server-held credentials supplied by client. The caller
// validates actual resource ownership before the tunnel is initialized.
// validateURL also verifies every returned relay endpoint before the SDK uses it.
func Open(ctx context.Context, name string, client *http.Client, authorize func(Codespace) error, validateURL func(string, string) error) (PortForwarder, error) {
	if client == nil || authorize == nil || validateURL == nil {
		return nil, fmt.Errorf("client and validators are required")
	}
	return open(ctx, name, client, authorize, validateURL)
}

func open(ctx context.Context, name string, client *http.Client, authorize func(Codespace) error, validateURL func(string, string) error) (PortForwarder, error) {
	f := &cmdutil.Factory{Config: func() (gh.Config, error) { return nil, fmt.Errorf("use default GitHub host") }, HttpClient: func() (*http.Client, error) { return client, nil }, ExternalHttpClient: func() (*http.Client, error) { return client, nil }}
	a := api.New(f)
	cs, err := a.GetCodespace(ctx, name, true)
	if err != nil {
		return nil, err
	}
	if err = authorize(Codespace{Name: cs.Name, Account: cs.Owner.Login, Owner: cs.Repository.Owner.Login, Repository: strings.TrimPrefix(cs.Repository.FullName, cs.Repository.Owner.Login+"/"), State: cs.State}); err != nil {
		return nil, err
	}
	if cs.State != api.CodespaceStateAvailable {
		return nil, fmt.Errorf("codespace is not available; start it before connecting")
	}
	if err = validateURL(cs.Connection.TunnelProperties.ServiceUri, "https"); err != nil {
		return nil, err
	}
	conn, err := connection.NewCodespaceConnectionWithRelayValidator(ctx, cs, client, func(raw string) error { return validateURL(raw, "wss") })
	if err != nil {
		return nil, err
	}
	return portforwarder.NewPortForwarder(ctx, conn)
}
