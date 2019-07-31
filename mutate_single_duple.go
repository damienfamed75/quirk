package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateSingleDupleNode(ctx context.Context, dg DgraphClient,
	node interface{}, uidMap map[string]UID, m *sync.Mutex) (bool, error) {

	return c.mutate(ctx, dg, node.(*DupleNode), uidMap, m)
}
