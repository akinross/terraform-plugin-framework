package float64planmodifier_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRequiresReplaceIfModifierPlanModifyFloat64(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.Float64Attribute{},
		},
	}

	nullPlan := tfsdk.Plan{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(context.Background()),
			nil,
		),
	}

	nullState := tfsdk.State{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(context.Background()),
			nil,
		),
	}

	testPlan := func(value types.Float64) tfsdk.Plan {
		tfValue, err := value.ToTerraformValue(context.Background())

		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Plan{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}

	testState := func(value types.Float64) tfsdk.State {
		tfValue, err := value.ToTerraformValue(context.Background())

		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.State{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}

	testCases := map[string]struct {
		request  planmodifier.Float64Request
		ifFunc   float64planmodifier.RequiresReplaceIfFunc
		expected *planmodifier.Float64Response
	}{
		"state-null": {
			// resource creation
			request: planmodifier.Float64Request{
				Plan:       testPlan(types.Float64Unknown()),
				PlanValue:  types.Float64Unknown(),
				State:      nullState,
				StateValue: types.Float64Null(),
			},
			ifFunc: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifier.RequiresReplaceIfFuncResponse) {
				resp.RequiresReplace = true // should never reach here
			},
			expected: &planmodifier.Float64Response{
				PlanValue:       types.Float64Unknown(),
				RequiresReplace: false,
			},
		},
		"plan-null": {
			// resource destroy
			request: planmodifier.Float64Request{
				Plan:       nullPlan,
				PlanValue:  types.Float64Null(),
				State:      testState(types.Float64Value(1.2)),
				StateValue: types.Float64Value(1.2),
			},
			ifFunc: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifier.RequiresReplaceIfFuncResponse) {
				resp.RequiresReplace = true // should never reach here
			},
			expected: &planmodifier.Float64Response{
				PlanValue:       types.Float64Null(),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-different-if-false": {
			request: planmodifier.Float64Request{
				Plan:       testPlan(types.Float64Value(2.4)),
				PlanValue:  types.Float64Value(2.4),
				State:      testState(types.Float64Value(1.2)),
				StateValue: types.Float64Value(1.2),
			},
			ifFunc: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifier.RequiresReplaceIfFuncResponse) {
				resp.RequiresReplace = false // no change
			},
			expected: &planmodifier.Float64Response{
				PlanValue:       types.Float64Value(2.4),
				RequiresReplace: false,
			},
		},
		"planvalue-statevalue-different-if-true": {
			request: planmodifier.Float64Request{
				Plan:       testPlan(types.Float64Value(2.4)),
				PlanValue:  types.Float64Value(2.4),
				State:      testState(types.Float64Value(1.2)),
				StateValue: types.Float64Value(1.2),
			},
			ifFunc: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifier.RequiresReplaceIfFuncResponse) {
				resp.RequiresReplace = true // should reach here
			},
			expected: &planmodifier.Float64Response{
				PlanValue:       types.Float64Value(2.4),
				RequiresReplace: true,
			},
		},
		"planvalue-statevalue-equal": {
			request: planmodifier.Float64Request{
				Plan:       testPlan(types.Float64Value(1.2)),
				PlanValue:  types.Float64Value(1.2),
				State:      testState(types.Float64Value(1.2)),
				StateValue: types.Float64Value(1.2),
			},
			ifFunc: func(ctx context.Context, req planmodifier.Float64Request, resp *float64planmodifier.RequiresReplaceIfFuncResponse) {
				resp.RequiresReplace = true // should never reach here
			},
			expected: &planmodifier.Float64Response{
				PlanValue:       types.Float64Value(1.2),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.Float64Response{
				PlanValue: testCase.request.PlanValue,
			}

			float64planmodifier.RequiresReplaceIf(testCase.ifFunc, "test", "test").PlanModifyFloat64(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
