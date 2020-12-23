package aws

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpecFromURL(t *testing.T) {
	cases := map[string]*Spec{
		"aws://profilename": {
			Profile: "profilename",
		},
		"aws://profilename?region=region": {
			Region: "region",
		},
		// Need to actually test this with env variables
		"aws://": nil,
	}
	for s, expected := range cases {
		u, err := url.Parse(s)
		if err != nil {
			t.Fatal(err)
		}
		got, err := SpecFromURL(u)
		if expected != nil {
			require.NoError(t, err)
			require.EqualValues(t, expected, got, s)
		} else {
			require.Error(t, err, s)
		}
	}
}
