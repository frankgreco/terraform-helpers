package overlap

import (
	"fmt"
	"net"
	"strings"

	"github.com/frankgreco/terraform-helpers/internal/utils"
)

func newCidrOrderedPair(cidr *net.IPNet) utils.OrderedPair {
	last := make([]byte, 4)
	for i := 0; i < 4; i++ {
		last[i] = cidr.IP[i] | (cidr.Mask[i] ^ 255)
	}
	return utils.NewOrderedPair(cidr.IP.Mask(cidr.Mask).To4().String(), net.IP(last).String())
}

func CIDR(encoded []string) error {
	if len(encoded) < 2 {
		return nil
	}

	var cidrs []utils.OrderedPair
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
			cidrs = append(cidrs, newCidrOrderedPair(ipNet))
		}
	}

	return utils.Overlaps(cidrs)
}
