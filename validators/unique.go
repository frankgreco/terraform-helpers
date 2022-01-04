package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	uniqueErr = "Unique constraint was violated for attribute %s."
)

type uniqueValidator struct {
	key string
}

// Unique returns an tfsdk.AttributeValidator that ensures
// all elements within a list or set are unique on the provided
// object attribute.
func Unique(key string) tfsdk.AttributeValidator {
	return uniqueValidator{
		key: key,
	}
}

func (u uniqueValidator) Description(context.Context) string {
	return fmt.Sprintf(uniqueErr, u.key)
}

func (u uniqueValidator) MarkdownDescription(context.Context) string {
	return fmt.Sprintf(uniqueErr, u.key)
}

func (u uniqueValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var items []types.Object

	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &items)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	log := map[interface{}]bool{}

	for _, item := range items {
		if item.Unknown || item.Null {
			continue
		}

		v, ok := item.Attrs[u.key]
		if !ok {
			continue // we'll rely on terraform enforcing `Required: true`
		}

		t, ok := item.AttrTypes[u.key]
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Internal Error with Terraform.",
				"There exists an attribute for which Terraform has not told us the type.",
			)
			return
		}

		var tmp interface{}

		switch t {
		case types.NumberType:
			var val int
			diags := tfsdk.ValueAs(ctx, v, &val)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
			tmp = val
		case types.StringType:
			var val string
			diags := tfsdk.ValueAs(ctx, v, &val)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
			tmp = val
		}

		if _, ok := log[tmp]; ok {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				fmt.Sprintf(uniqueErr, u.key),
				fmt.Sprintf("More than one item exists with %s=%v", u.key, tmp),
			)
			return
		}
		log[tmp] = true

	}
}
