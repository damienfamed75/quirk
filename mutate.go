package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutate(ctx context.Context, dg DgraphClient,
	d *DupleNode, uidMap map[string]string, m sync.Locker) (bool, error) {

	res := c.tryUpsert(ctx, dg.NewTxn(), d)

	if res.err != nil {
		return res.new, res.err
	}
	if res.new {
		// Note: The reason this isn't in the worker is because
		// when calling to add a single node then this map will
		// not be updated.
		m.Lock()
		uidMap[res.identifier] = res.uid
		m.Unlock()
	}

	return res.new, nil
}
