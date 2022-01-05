package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	conflictsWithErr         = "There was a conflict detected."
	conflictsWithDescription = "Ensures that the specificed attributes at the same level are not set (either null or unknown)."
)

type conflictsWithValidator struct {
	conflicts []string
}

// ConflictsWith ensures that the specificed attributes at the same level are not set (either null or unknown).
func ConflictsWith(attributes ...string) tfsdk.AttributeValidator {
	return conflictsWithValidator{
		conflicts: attributes,
	}
}

// Description describes this validator.
func (v conflictsWithValidator) Description(context.Context) string {
	return conflictsWithDescription
}

// MarkdownDescription describes this validator.
func (v conflictsWithValidator) MarkdownDescription(context.Context) string {
	return conflictsWithDescription
}

// Validate performs validation on an attribute.
func (v conflictsWithValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if len(v.conflicts) == 0 {
		return
	}

	this, err := toValue(ctx, req.AttributeConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			conflictsWithErr,
			"The validator had an internal error: "+err.Error(),
		)
		return
	}

	// We don't need to do any validation if the value isn't "set".
	if !this.IsFullyKnown() || this.IsNull() {
		return
	}

	var parent types.Object
	{
		// I think this is almost always guaranteed to be an Object?
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.AttributePath.WithoutLastStep(), &parent)...)
		if resp.Diagnostics.HasError() || parent.Null || parent.Unknown {
			return
		}
	}

	conflicts := []string{}
	for _, conflict := range v.conflicts {
		// I believe the only way this wouldn't be true
		// is if they pass in an unknown attribute.
		attrValue, ok := parent.Attrs[conflict]
		if !ok {
			continue
		}

		data, err := toValue(ctx, attrValue)
		if err != nil {
			resp.Diagnostics.AddError(
				conflictsWithErr,
				"The validator had an internal error: "+err.Error(),
			)
			return
		}

		// Check if the attribute is "actually" set.
		if data.IsFullyKnown() && !data.IsNull() {
			conflicts = append(conflicts, conflict)
		}
	}

	if len(conflicts) > 0 {
		resp.Diagnostics.AddError(
			conflictsWithErr,
			req.AttributePath.String()+" conficts with "+strings.Join(conflicts, ", ")+".",
		)
	}
}
