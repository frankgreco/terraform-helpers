package validators

import (
	"context"
	"strconv"

	"github.com/frankgreco/terraform-helpers/internal/overlap"
	"github.com/frankgreco/terraform-helpers/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	noOverlapErr         = "There was an overlap detected."
	noOverlapDescription = "Ensures that no items overlap with any other in the list."
)

type noOverlapValidator struct{}

// NoOverlap ensures that no elements overlap with any other in the list.
func NoOverlap() tfsdk.AttributeValidator {
	return noOverlapValidator{}
}

// Description describes this validator.
func (v noOverlapValidator) Description(context.Context) string {
	return noOverlapDescription
}

// MarkdownDescription describes this validator.
func (v noOverlapValidator) MarkdownDescription(context.Context) string {
	return noOverlapDescription
}

// Validate performs validation on an attribute.
func (v noOverlapValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var list types.List
	{
		resp.Diagnostics.Append(tfsdk.ValueAs(ctx, req.AttributeConfig, &list)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if list.Unknown || list.Null {
			return
		}
	}

	// Object
	if (tftypes.Object{}).Is(list.Type(ctx).(types.ListType).ElementType().TerraformType(ctx)) {
		resp.Diagnostics.Append(validateObjectList(ctx, list)...)
		return
	}

	// Number
	switch list.Type(ctx).(types.ListType).ElementType() {
	case types.NumberType:
		resp.Diagnostics.Append(validateNumberList(ctx, list)...)
	default:
		resp.Diagnostics.AddError(
			noOverlapErr,
			"Unsupported list element type.",
		)
		return
	}
}

func validateObjectList(ctx context.Context, list types.List) (diags diag.Diagnostics) {
	var items []types.Object
	{
		diags.Append(list.ElementsAs(ctx, &items, false)...)
		if diags.HasError() {
			return
		}
	}

	var pairs []utils.OrderedPair
	for _, item := range items {
		aux := struct {
			From int `tfsdk:"from"`
			To   int `tfsdk:"to"`
		}{}
		diags.Append(item.As(ctx, &aux, types.ObjectAsOptions{})...)
		if diags.HasError() {
			return
		}
		pairs = append(pairs, utils.NewOrderedPair(strconv.Itoa(aux.From), strconv.Itoa(aux.To)))
	}

	if err := overlap.OrderedPairs(pairs); err != nil {
		diags.AddError(noOverlapErr, err.Error())
		return
	}
	return
}

func validateNumberList(ctx context.Context, list types.List) (diags diag.Diagnostics) {
	var encoded []int
	{
		diags.Append(list.ElementsAs(ctx, &encoded, false)...)
		if diags.HasError() {
			return
		}
	}

	if err := overlap.IntSlice(encoded); err != nil {
		diags.AddError(
			noOverlapErr,
			err.Error(),
		)
		return
	}
	return
}
