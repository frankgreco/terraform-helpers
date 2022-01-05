package cidr

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOverlaps(t *testing.T) {
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
			name: "overlapping happy path",
			cidrs: []string{
				"192.168.1.0/24",
				"192.168.2.0/16",
			},
			err: "The IPs between 192.168.1.0 and 192.168.255.255 are supplied by more than one range.",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cidrs := []*net.IPNet{}
			for _, cidr := range test.cidrs {
				_, ipNet, err := net.ParseCIDR(cidr)
				require.NoError(t, err, test.name)
				require.NotNil(t, ipNet.IP.To4(), test.name)
				cidrs = append(cidrs, ipNet)
			}
			err := Overlaps(cidrs)
			if test.err == "" {
				require.NoError(t, err, test.name)
			} else {
				require.Equal(t, test.err, err.Error(), test.name)
			}
		})
	}
}

func TestNewCidrRange(t *testing.T) {
	for _, test := range []struct {
		name     string
		cidr     string
		expected cidrRange
	}{
		{
			name: "network boundary /24",
			cidr: "192.168.2.0/24",
			expected: cidrRange{
				first: []byte{192, 168, 2, 0},
				last:  []byte{192, 168, 2, 255},
			},
		},
		{
			name: "non-network boundary /24",
			cidr: "192.168.2.10/24",
			expected: cidrRange{
				first: []byte{192, 168, 2, 0},
				last:  []byte{192, 168, 2, 255},
			},
		},
		{
			name: "non-network boundary /16",
			cidr: "192.168.2.10/16",
			expected: cidrRange{
				first: []byte{192, 168, 0, 0},
				last:  []byte{192, 168, 255, 255},
			},
		},
		{
			name: "network boundary /8",
			cidr: "192.0.0.0/8",
			expected: cidrRange{
				first: []byte{192, 0, 0, 0},
				last:  []byte{192, 255, 255, 255},
			},
		},
		{
			name: "/32",
			cidr: "192.168.2.42/32",
			expected: cidrRange{
				first: []byte{192, 168, 2, 42},
				last:  []byte{192, 168, 2, 42},
			},
		},
		{
			name: "non octet boundary /17",
			cidr: "192.168.86.42/17",
			expected: cidrRange{
				first: []byte{192, 168, 0, 0},
				last:  []byte{192, 168, 127, 255},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, ipNet, err := net.ParseCIDR(test.cidr)
			require.NoError(t, err, test.name)
			require.Equal(t, test.expected, newCidrRange(ipNet), test.name)
		})
	}
}
