package validators

import (
	"context"
	"net"
	"strings"

	"github.com/frankgreco/terraform-helpers/internal/cidr"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	noOverlappingCIDRsErr         = "There was an overlap detected."
	noOverlappingCIDRsDescription = "Ensures that no CIDRs overlap with any other in the list."
)

type noOverlappingCIDRsValidator struct{}

// NoOverlappingCIDRs ensures that no CIDRs overlap with any other in the list.
func NoOverlappingCIDRs() tfsdk.AttributeValidator {
	return noOverlappingCIDRsValidator{}
}

// Description describes this validator.
func (v noOverlappingCIDRsValidator) Description(context.Context) string {
	return noOverlappingCIDRsDescription
}

// MarkdownDescription describes this validator.
func (v noOverlappingCIDRsValidator) MarkdownDescription(context.Context) string {
	return noOverlappingCIDRsDescription
}

// Validate performs validation on an attribute.
func (v noOverlappingCIDRsValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {

	var list types.List
	{
		resp.Diagnostics.Append(tfsdk.ValueAs(ctx, req.AttributeConfig, &list)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if list.Unknown || list.Null {
			return
		}
	}

	var encoded []string
	{
		resp.Diagnostics.Append(list.ElementsAs(ctx, &encoded, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(encoded) < 2 {
			return
		}
	}

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
				resp.Diagnostics.AddError(
					noOverlappingCIDRsErr,
					"Invalid CIDR: "+err.Error(),
				)
				return
			}
			if ipNet.IP.To4() == nil {
				resp.Diagnostics.AddError(
					noOverlappingCIDRsErr,
					"Invalid CIDR: not IPv4",
				)
				return
			}
			cidrs = append(cidrs, ipNet)
		}
	}

	if err := cidr.Overlaps(cidrs); err != nil {
		resp.Diagnostics.AddError(
			noOverlappingCIDRsErr,
			err.Error(),
		)
		return
	}

}
