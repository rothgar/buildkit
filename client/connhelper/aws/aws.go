// Package aws provides connhelper for aws://<profile>
package aws

import (
	"context"
	"net"
	"net/url"
	"os"

	"github.com/docker/cli/cli/connhelper/commandconn"
	"github.com/moby/buildkit/client/connhelper"
	"github.com/pkg/errors"
)

func init() {
	connhelper.Register("aws", Helper)
}

// Helper returns helper for connecting to an AWS account.
// Requires BuildKit v0.5.0 or later in the pod.
func Helper(u *url.URL) (*connhelper.ConnectionHelper, error) {
	sp, err := SpecFromURL(u)
	if err != nil {
		return nil, err
	}
	return &connhelper.ConnectionHelper{
		ContextDialer: func(ctx context.Context, addr string) (net.Conn, error) {
			// using background context because context remains active for the duration of the process, after dial has completed
			return commandconn.New(context.Background(), "aws", "--profile="+sp.Profile, "--region="+sp.Region,
				"ssm", "start-session", "--target="+sp.Instance, "--document-name=buildkit")
		},
	}, nil
}

// Spec
type Spec struct {
	Instance string
	Profile  string
	Region   string
}

// SpecFromURL creates Spec from URL.
// URL is like aws://instance?profile=<profile>,region=<region>
// Only <instance> part is mandatory.
func SpecFromURL(u *url.URL) (*Spec, error) {
	q := u.Query()
	var awsProfileSet bool
	sp := Spec{
		Instance: u.Hostname(),
		Profile:  q.Get("profile"),
		Region:   q.Get("region"),
	}
	if sp.Instance == "" {
		return nil, errors.New("Need to provide an instance id")
	}
	if sp.Profile == "" {
		var awsDefaultProfileSet bool
		sp.Profile, awsDefaultProfileSet = os.LookupEnv("AWS_DEFAULT_PROFILE")
		if !awsDefaultProfileSet {
			sp.Profile, awsProfileSet = os.LookupEnv("AWS_PROFILE")
			if !awsProfileSet {
				return nil, errors.New("AWS profile not set")
			}
		}
	}
	if sp.Region == "" {
		var awsDefaultRegionSet bool
		sp.Region, awsDefaultRegionSet = os.LookupEnv("AWS_DEFAULT_REGION")
		if !awsDefaultRegionSet {
			sp.Region, _ = os.LookupEnv("AWS_REGION")
		}
	}
	return &sp, nil
}
