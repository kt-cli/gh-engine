// Package transport defines an optional managed port forwarding boundary.
package transport

import (
	"context"
	"github.com/microsoft/dev-tunnels/go/tunnels"
	"io"
	"net"
)

// ForwardPortOpts selects a remote port and its forwarding behavior.
type ForwardPortOpts struct {
	Port       int
	Internal   bool
	KeepAlive  bool
	Visibility string
}

// PortForwarder preserves the command core's port forwarding operations.
type PortForwarder interface {
	ForwardPortToListener(context.Context, ForwardPortOpts, *net.TCPListener) error
	ForwardPort(context.Context, ForwardPortOpts) error
	ConnectToForwardedPort(context.Context, io.ReadWriteCloser, ForwardPortOpts) error
	ListPorts(context.Context) ([]*tunnels.TunnelPort, error)
	UpdatePortVisibility(context.Context, int, string) error
	KeepAlive(string)
	GetKeepAliveReason() string
	Close() error
}

// Connector may be implemented by an injected external HTTP transport.
// Its implementation owns connection credentials and may return a remote forwarder.
type Connector interface {
	ConnectCodespace(context.Context, string) (PortForwarder, error)
}
