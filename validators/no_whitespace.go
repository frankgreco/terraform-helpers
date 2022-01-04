package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	noWhitespaceValidatorErr = "string must not contain whitespace"
)

type noWhitespaceValidator struct{}

func NoWhitespace() tfsdk.AttributeValidator {
	return noWhitespaceValidator{}
}

func (v noWhitespaceValidator) Description(context.Context) string {
	return noWhitespaceValidatorErr
}

func (v noWhitespaceValidator) MarkdownDescription(context.Context) string {
	return noWhitespaceValidatorErr
}

func (v noWhitespaceValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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

	if strings.Contains(str.Value, " ") {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Content",
			noWhitespaceValidatorErr,
		)
	}
}
