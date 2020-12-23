// Package ssh provides connhelper for ssh://<hostname>
package ssh

import (
	"context"
	"net"
	"net/url"

	"github.com/docker/cli/cli/connhelper/commandconn"
	"github.com/moby/buildkit/client/connhelper"
	"github.com/pkg/errors"
)

func init() {
	connhelper.Register("ssh", Helper)
}

// Helper returns helper for connecting to an ssh backend.
// Requires BuildKit v0.5.0 or later on the host.
func Helper(u *url.URL) (*connhelper.ConnectionHelper, error) {
	sp, err := SpecFromURL(u)
	if err != nil {
		return nil, err
	}
	return &connhelper.ConnectionHelper{
		ContextDialer: func(ctx context.Context, addr string) (net.Conn, error) {
			// using background context because context remains active for the duration of the process, after dial has completed
			return commandconn.New(context.Background(), "ssh", sp.Hostname, "--", "buildctl", "dial-stdio")
		},
	}, nil
}

// Spec
type Spec struct {
	Hostname string
}

// SpecFromURL creates Spec from URL.
// URL is like ssh://hostname
// Only <hostname> part is mandatory.
func SpecFromURL(u *url.URL) (*Spec, error) {
	sp := Spec{
		Hostname: u.Hostname(),
	}
	if sp.Hostname == "" {
		return nil, errors.New("Need to provide a hostname")
	}
	return &sp, nil
}
