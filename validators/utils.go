package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func toValue(ctx context.Context, in attr.Value) (tftypes.Value, error) {
	data, err := in.ToTerraformValue(ctx)
	if err != nil {
		return tftypes.Value{}, err
	}
	return tftypes.NewValue(in.Type(ctx).TerraformType(ctx), data), nil
}
