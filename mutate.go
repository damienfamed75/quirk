package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutate(ctx context.Context, dg DgraphClient,
	d *DupleNode, uidMap map[string]UID, m sync.Locker) (bool, error) {

	res := c.tryUpsert(ctx, dg.NewTxn(), d)
	if res.err != nil {
		return res.new, res.err
	}

	m.Lock()
	uidMap[res.identifier] = UID{uid: res.uid, isNew: res.new}
	m.Unlock()

	return res.new, nil
}
