package validators

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	cidrErr = "value must be a valid cidr"
)

type cidrValidator struct{}

func Cidr() tfsdk.AttributeValidator {
	return cidrValidator{}
}

func (v cidrValidator) Description(context.Context) string {
	return cidrErr
}

func (v cidrValidator) MarkdownDescription(context.Context) string {
	return cidrErr
}

func (v cidrValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var str types.String
	{
		diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &str)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	if str.Unknown || str.Null {
		return
	}

	if _, _, err := net.ParseCIDR(str.Value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Content",
			cidrErr,
		)
	}
}
