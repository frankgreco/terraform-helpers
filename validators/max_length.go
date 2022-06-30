package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	maxLengthErr         = "String must be at most %d characters long."
	maxLengthDescription = "Ensure that the attribute value has a maximum length."
)

type maxLengthValidator struct {
	length int
}

func MaxLength(length int) tfsdk.AttributeValidator {
	return maxLengthValidator{
		length: length,
	}
}

func (v maxLengthValidator) Description(context.Context) string {
	return maxLengthDescription
}

func (v maxLengthValidator) MarkdownDescription(context.Context) string {
	return maxLengthDescription
}

func (v maxLengthValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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

	if len(str.Value) > v.length {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Length",
			fmt.Sprintf(maxLengthErr, v.length),
		)
	}
}
