package overlap

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCIDR(t *testing.T) {
	for _, test := range []struct {
		name  string
		cidrs []string
		err   string
	}{
		{
			name: "no overlapping happy path",
			cidrs: []string{
				"192.168.1.0/24",
				"192.168.2.0/24",
				"192.168.3.0/24",
			},
		},
		{
			name: "no mask",
			cidrs: []string{
				"192.168.1.0/24",
				"192.168.2.0/24",
				"192.168.3.0",
			},
		},
		{
			name: "no mask overlap",
			cidrs: []string{
				"192.168.1.0/24",
				"192.168.2.0/24",
				"192.168.3.0",
				"192.168.3.0/32",
			},
			err: "The element 192.168.3.0 is supplied by more than one range.",
		},
		{
			name: "overlapping happy path",
			cidrs: []string{
				"192.168.1.0/24",
				"192.168.2.0/16",
			},
			err: "The elements between 192.168.1.0 and 192.168.255.255 are supplied by more than one range.",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := CIDR(test.cidrs); test.err == "" {
				require.NoError(t, err, test.name)
			} else {
				require.NotNil(t, err, test.name)
				require.Equal(t, test.err, err.Error(), test.name)
			}
		})
	}
}

func TestNewCidrRange(t *testing.T) {
	for _, test := range []struct {
		name   string
		cidr   string
		first  string
		second string
	}{
		{
			name:   "network boundary /24",
			cidr:   "192.168.2.0/24",
			first:  "192.168.2.0",
			second: "192.168.2.255",
		},
		{
			name:   "non-network boundary /24",
			cidr:   "192.168.2.10/24",
			first:  "192.168.2.0",
			second: "192.168.2.255",
		},
		{
			name:   "non-network boundary /16",
			cidr:   "192.168.2.10/16",
			first:  "192.168.0.0",
			second: "192.168.255.255",
		},
		{
			name:   "network boundary /8",
			cidr:   "192.0.0.0/8",
			first:  "192.0.0.0",
			second: "192.255.255.255",
		},
		{
			name:   "/32",
			cidr:   "192.168.2.42/32",
			first:  "192.168.2.42",
			second: "192.255.2.42",
		},
		{
			name:   "non octet boundary /17",
			cidr:   "192.168.86.42/17",
			first:  "192.168.0.0",
			second: "192.255.127.255",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, ipNet, err := net.ParseCIDR(test.cidr)
			require.NoError(t, err, test.name)
			op := newCidrOrderedPair(ipNet)
			require.Equal(t, test.first, op.First(), test.name)
			require.Equal(t, test.second, op.Last(), test.name)
		})
	}
}
