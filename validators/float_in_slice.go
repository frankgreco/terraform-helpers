package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	floatInSliceErr         = "number must be one of %s"
	floatInSliceDescription = "Ensure that the attribute value is one of the provided values."
)

type floatInSliceValidator struct {
	values []float64
}

func FloatInSlice(values ...float64) tfsdk.AttributeValidator {
	return floatInSliceValidator{
		values: values,
	}
}

func (v floatInSliceValidator) Description(context.Context) string {
	return floatInSliceDescription
}

func (v floatInSliceValidator) MarkdownDescription(context.Context) string {
	return floatInSliceDescription
}

func (v floatInSliceValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var number types.Float64
	{
		diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &number)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	if number.Unknown || number.Null {
		return
	}

	for _, val := range v.values {
		if val == number.Value {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		"Invalid Number",
		fmt.Sprintf(floatInSliceErr, strings.Replace(fmt.Sprint(v.values), " ", ", ", -1)),
	)
}
