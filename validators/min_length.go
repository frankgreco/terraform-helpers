package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	lengthValidatorErr = "String must be at least %d characters long."
)

type lengthValidator struct {
	length int
}

func MinLength(length int) tfsdk.AttributeValidator {
	return lengthValidator{
		length: length,
	}
}

func (v lengthValidator) Description(context.Context) string {
	return fmt.Sprintf(lengthValidatorErr, v.length)
}

func (v lengthValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf(lengthValidatorErr, v.length)
}

func (v lengthValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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
			fmt.Sprintf(lengthValidatorErr, v.length),
		)
	}
}
