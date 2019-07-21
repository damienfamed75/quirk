package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateSingleStruct(ctx context.Context, dg DgraphClient,
	d interface{}, uidMap map[string]string, m *sync.Mutex) (bool, error) {
	// Use reflect to package the predicate and values in slices.
	predVals := c.reflectMaps(d)

	return c.mutate(ctx, dg, predVals, uidMap, m)
}
