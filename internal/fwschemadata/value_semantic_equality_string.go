package fwschemadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValueSemanticEqualityString performs string type semantic equality.
func ValueSemanticEqualityString(ctx context.Context, req ValueSemanticEqualityRequest, resp *ValueSemanticEqualityResponse) {
	priorValuable, ok := req.PriorValue.(basetypes.StringValuableWithSemanticEquals)

	// No changes required if the interface is not implemented.
	if !ok {
		return
	}

	proposedNewValuable, ok := req.ProposedNewValue.(basetypes.StringValuableWithSemanticEquals)

	// No changes required if the interface is not implemented.
	if !ok {
		return
	}

	logging.FrameworkTrace(
		ctx,
		"Calling provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	usePriorValue, diags := proposedNewValuable.StringSemanticEquals(ctx, priorValuable)

	logging.FrameworkTrace(
		ctx,
		"Called provider defined type-based SemanticEquals",
		map[string]interface{}{
			logging.KeyValueType: proposedNewValuable.String(),
		},
	)

	resp.Diagnostics.Append(diags...)

	if !usePriorValue {
		return
	}

	resp.NewValue = priorValuable
}
