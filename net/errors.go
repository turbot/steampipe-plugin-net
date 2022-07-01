package net

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func shouldRetryError() plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		if err != nil {
			panic(d.Table.Name)
			return true
		}
		return false
	}
}
