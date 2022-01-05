package cidr

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

type cidrRange struct {
	first []byte
	last  []byte
}

func newCidrRange(cidr *net.IPNet) cidrRange {
	out := cidrRange{
		first: cidr.IP.Mask(cidr.Mask).To4(),
		last:  make([]byte, 4),
	}

	for i := 0; i < 4; i++ {
		out.last[i] = cidr.IP[i] | (cidr.Mask[i] ^ 255)
	}

	return out
}

func Overlaps(encoded []string) error {
	var cidrs []*net.IPNet
	{
		for _, cidr := range encoded {
			// If the cidr contains a '/' anywhere else, It's already malformed.
			// Making it more malformed won't be an issue.
			if !strings.Contains(cidr, "/") {
				cidr = cidr + "/32"
			}

			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				return err
			}
			if ipNet.IP.To4() == nil {
				return fmt.Errorf("Invalid CIDR: %s is not IPv4", ipNet.IP.String())
			}
			cidrs = append(cidrs, ipNet)
		}
	}

	return overlaps(cidrs)
}

func overlaps(cidrs []*net.IPNet) error {
	if len(cidrs) < 2 {
		return nil
	}

	var ranges []cidrRange
	for _, cidr := range cidrs {
		if cidr == nil {
			continue
		}
		ranges = append(ranges, newCidrRange(cidr))
	}
	sort.Slice(ranges, func(i, j int) bool {
		return net.IP(ranges[i].first).String() < net.IP(ranges[j].first).String()
	})

	for i := 1; i < len(ranges); i++ {
		if l, f := net.IP(ranges[i-1].last).String(), net.IP(ranges[i].first).String(); l >= f {
			return fmt.Errorf("The IPs between %s and %s are supplied by more than one range.", f, l)
		}
	}

	return nil
}
