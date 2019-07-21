package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateSingleDupleNode(ctx context.Context, dg DgraphClient,
	node *DupleNode, uidMap map[string]string, m sync.Locker) (bool, error) {

	return c.mutate(ctx, dg, node, uidMap, m)
}
