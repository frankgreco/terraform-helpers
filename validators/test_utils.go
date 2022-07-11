package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name      string
	validator tfsdk.AttributeValidator
	request   tfsdk.ValidateAttributeRequest
	err       bool
}

func (tc testCase) run(t *testing.T) {
	response := tfsdk.ValidateAttributeResponse{
		Diagnostics: diag.Diagnostics{},
	}

	tc.validator.Validate(context.Background(), tc.request, &response)

	if hasError := response.Diagnostics.HasError(); tc.err {
		assert.True(t, hasError, "diagnostic did not contain expected error")
	} else {
		assert.False(t, hasError, "diagnostic contained an unexpected error")
	}

	// TODO: assert the acutal diagnostics
}
