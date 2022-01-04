package validators

import (
	"context"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	conflictsWithErr = ""
)

type conflictsWithValidator struct{}

func ConflictsWith(attributes ...string) tfsdk.AttributeValidator {
	return conflictsWithValidator{}
}

func (v conflictsWithValidator) Description(context.Context) string {
	return conflictsWithErr
}

func (v conflictsWithValidator) MarkdownDescription(context.Context) string {
	return conflictsWithErr
}

func (v conflictsWithValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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
