package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateStringMap(ctx context.Context, dg DgraphClient,
	d map[string]string, uidMap map[string]UID, m sync.Locker) (bool, error) {
	// Convert out map[string]string to usable predicate and value data.
	predVals := c.mapToPredValPairs(d)

	return c.mutate(ctx, dg, predVals, uidMap, m)
}

func (c *Client) mutateDynamicMap(ctx context.Context, dg DgraphClient,
	d map[string]interface{}, uidMap map[string]UID, m sync.Locker) (bool, error) {
	// Convert out map[string]interface{} to usable predicate and value data.
	predVals := c.dynamicMapToPredValPairs(d)

	return c.mutate(ctx, dg, predVals, uidMap, m)
}
