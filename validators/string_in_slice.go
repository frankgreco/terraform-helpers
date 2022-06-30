package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	stringInSliceErr         = "string must be one of [%s]"
	stringInSliceDescription = "Ensure that the attribute value is one of the provided values."
)

type stringInSliceValidator struct {
	caseSensitive bool
	values        []string
}

func StringInSlice(caseSensitive bool, values ...string) tfsdk.AttributeValidator {
	return stringInSliceValidator{
		caseSensitive: caseSensitive,
		values:        values,
	}
}

func (v stringInSliceValidator) Description(context.Context) string {
	return stringInSliceDescription
}

func (v stringInSliceValidator) MarkdownDescription(context.Context) string {
	return stringInSliceDescription
}

func (v stringInSliceValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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

	for _, val := range v.values {
		if val == str.Value {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid String",
		fmt.Sprintf(stringInSliceErr, strings.Join(v.values, ", ")),
	)
}
