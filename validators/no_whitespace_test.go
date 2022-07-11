package validators

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNoWhitespace(t *testing.T) {
	for _, test := range []testCase{
		{
			name:      "pass",
			validator: NoWhitespace(),
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{
					Value: "no_whitespace",
				},
			},
		},
		{
			name:      "fail",
			validator: NoWhitespace(),
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{
					Value: "no whitespace",
				},
			},
			err: true,
		},
		{
			name:      "null",
			validator: NoWhitespace(),
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{
					Null: true,
				},
			},
		},
		{
			name:      "unknown",
			validator: NoWhitespace(),
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.String{
					Unknown: true,
				},
			},
		},
		{
			name:      "wrong type",
			validator: NoWhitespace(),
			request: tfsdk.ValidateAttributeRequest{
				AttributeConfig: types.Bool{
					Value: true,
				},
			},
			err: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.run(t)
		})
	}
}
