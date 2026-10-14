// Package ktengine exposes the upstream command tree with a credential-free
// configuration and injected HTTP clients. Authorization belongs to the transport.
package ktengine

import (
	"fmt"
	"github.com/cli/cli/v2/internal/config"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/internal/telemetry"
	"github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/cli/cli/v2/pkg/cmd/root"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

// Version identifies the immutable upstream release used by this adapter.
const Version = "2.102.0"

// New constructs upstream commands using the provided streams and transport.
// The caller must reject or remove commands that export credentials or load
// arbitrary executables before executing the returned command tree.
func New(ios *iostreams.IOStreams, client *http.Client, external *http.Client, executable string) (*cobra.Command, error) {
	if client == nil || client.Transport == nil || external == nil || external.Transport == nil {
		return nil, fmt.Errorf("explicit HTTP transports are required")
	}
	cfg := config.NewTransportConfig()
	t := &telemetry.NoOpService{}
	f := factory.New(Version, "kt-gh", func() (gh.Config, error) { return cfg, nil }, ios, executable, t)
	f.APITransport = client.Transport
	f.HttpClient = func() (*http.Client, error) { return client, nil }
	f.PlainHttpClient = f.HttpClient
	f.ExternalHttpClient = func() (*http.Client, error) { return external, nil }
	// Extensions are not loaded from the operator's filesystem at construction.
	f.ExtensionManager = noExtensions{}
	cmd, err := root.NewCmdRoot(f, t, Version, "")
	if err != nil {
		return nil, err
	}
	cmd.Use = "kt-gh <command> <subcommand> [flags]"
	return cmd, nil
}
