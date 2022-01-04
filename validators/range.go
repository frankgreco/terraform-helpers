package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	rangeErr = "value must be between %v and %v"
)

type rangeValidator struct {
	from, to interface{}
}

func Range(from, to interface{}) tfsdk.AttributeValidator {
	return rangeValidator{
		from: from,
		to:   to,
	}
}

func (v rangeValidator) Description(context.Context) string {
	return fmt.Sprintf(rangeErr, v.from, v.to)
}

func (v rangeValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf(rangeErr, v.from, v.to)
}

func (v rangeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	switch req.AttributeConfig.Type(ctx) {
	case types.NumberType:
		var val types.Number
		diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &val)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		if val.Unknown || val.Null {
			return
		}

		f, ok := v.from.(float64)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Invalid From Value",
				"This validator was initialized with an incorrect type for from",
			)
			return
		}

		t, ok := v.to.(float64)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Invalid From Value",
				"This validator was initialized with an incorrect type for to",
			)
			return
		}

		x, _ := val.Value.Float64()

		if x < f || x > t {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Invalid Value",
				fmt.Sprintf(rangeErr, v.from, v.to),
			)
		}
	}
}
