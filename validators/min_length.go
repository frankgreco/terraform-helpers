package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	minLengthErr         = "String must be at least %d characters long."
	minLengthDescription = "Ensure that the attribute value has a minimum length."
)

type minLengthValidator struct {
	length int
}

func MinLength(length int) tfsdk.AttributeValidator {
	return minLengthValidator{
		length: length,
	}
}

func (v minLengthValidator) Description(context.Context) string {
	return minLengthDescription
}

func (v minLengthValidator) MarkdownDescription(context.Context) string {
	return minLengthDescription
}

func (v minLengthValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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

	if len(str.Value) < v.length {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Length",
			fmt.Sprintf(minLengthErr, v.length),
		)
	}
}
