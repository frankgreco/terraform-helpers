package cidr

import (
	"fmt"
	"net"
	"sort"
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

func Overlaps(cidrs []*net.IPNet) error {
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
