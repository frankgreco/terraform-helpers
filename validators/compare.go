package validators

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	compareErr         = "The comparison failed."
	compareDescription = "Ensures that the specificed attributes at the same level are not set (either null or unknown)."
)

type Comparator int

const (
	ComparatorUnknown Comparator = iota
	ComparatorLessThan
	ComparatorGreaterThan
	ComparatorEqual
	ComparatorLessThanEqual
	ComparatorGreaterThanEqual
	ComparatorNot
)

type compareValidator struct {
	comparator Comparator
	attribute  string
}

// Compare
func Compare(comparator Comparator, attribute string) tfsdk.AttributeValidator {
	return compareValidator{
		comparator: comparator,
		attribute:  attribute,
	}
}

// Description describes this validator.
func (v compareValidator) Description(context.Context) string {
	return compareDescription
}

// MarkdownDescription describes this validator.
func (v compareValidator) MarkdownDescription(context.Context) string {
	return compareDescription
}

// Validate performs validation on an attribute.
func (v compareValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if v.comparator == ComparatorUnknown {
		resp.Diagnostics.AddError(
			compareErr,
			"Unknown comparator",
		)
		return
	}

	this, err := toValue(ctx, req.AttributeConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			compareErr,
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

	// I believe the only way this wouldn't be true
	// is if they pass in an unknown attribute.
	attrValue, ok := parent.Attrs[v.attribute]
	if !ok {
		return
	}

	data, err := toValue(ctx, attrValue)
	if err != nil {
		resp.Diagnostics.AddError(
			compareErr,
			"The validator had an internal error: "+err.Error(),
		)
		return
	}

	// Check if the attribute is "actually" set.
	if !data.IsFullyKnown() || data.IsNull() {
		return
	}

	if err := compare(this, data, v.comparator); err != nil {
		resp.Diagnostics.AddError(compareErr, err.Error())
		return
	}
}

func compare(left, right tftypes.Value, comparator Comparator) error {
	if !left.Type().Is(right.Type()) {
		return errors.New("The type of both operands must match")
	}

	// TOOD: Add other types as needed.
	if left.Type().Is(tftypes.Number) {
		var leftOp, rightOp *big.Float
		if err := left.As(&leftOp); err != nil {
			return errors.New("The validator had an internal error: " + err.Error())
		}
		if err := right.As(&rightOp); err != nil {
			return errors.New("The validator had an internal error: " + err.Error())
		}
		return compareNumber(leftOp, rightOp, comparator)
	}
	return errors.New("Unsupported type. Only string and number are supported.")
}

func compareNumber(l, r *big.Float, comparator Comparator) error {
	left, _ := l.Float64()
	right, _ := r.Float64()

	// Using the following pattern for readability.
	switch comparator {
	case ComparatorLessThan:
		if left < right {
		} else {
			return fmt.Errorf("%f is not less than %f", left, right)
		}
	case ComparatorGreaterThan:
		if left > right {
		} else {
			return fmt.Errorf("%f is not greater than %f", left, right)
		}
	case ComparatorEqual:
		if left == right {
		} else {
			return fmt.Errorf("%f is not equal to %f", left, right)
		}
	case ComparatorLessThanEqual:
		if left <= right {
		} else {
			return fmt.Errorf("%f is not less than or equal to %f", left, right)
		}
	case ComparatorGreaterThanEqual:
		if left >= right {
		} else {
			return fmt.Errorf("%f is not greater than or equal to %f", left, right)
		}
	case ComparatorNot:
		if left != right {
		} else {
			return fmt.Errorf("%f is not (not equal) to %f", left, right)
		}
	default:
		return errors.New("Unknown comparator")
	}
	return nil
}
