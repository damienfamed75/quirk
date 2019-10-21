package quirk

import (
	"context"
	"fmt"
	"sync"

	"github.com/dgraph-io/dgo/v2"
)

// mutate is used to upsert a single node using Quirk. All single mutate
// functions call this one to do the upsert.
func (c *Client) mutate(ctx context.Context, dg *dgo.Dgraph,
	d *DupleNode, uidMap map[string]UID, m sync.Locker) (bool, error) {

	res := c.tryUpsert(ctx, dg.NewTxn(), d)
	if res.err != nil {
		return res.new, fmt.Errorf("mutate: %w", res.err)
	}

	m.Lock()
	uidMap[res.identifier] = UID{uid: res.uid, isNew: res.new}
	m.Unlock()

	return res.new, nil
}

// mutateSingleStruct uses reflect to create a DupleNode struct out of the given
// structure. Once the DupleNode is made then it is passed to mutate to be upserted.
func (c *Client) mutateSingleStruct(ctx context.Context, dg *dgo.Dgraph,
	d interface{}, uidMap map[string]UID, m *sync.Mutex) (bool, error) {
	// Use reflect to package the predicate and values in slices.
	predVals := c.reflectMaps(d)

	n, err := c.mutate(ctx, dg, predVals, uidMap, m)
	if err != nil {
		return n, fmt.Errorf("mutateSingleStruct: %w", err)
	}

	return n, nil
}

// mutateSingleDupleNode passed the DupleNode given to mutate to be upserted.
func (c *Client) mutateSingleDupleNode(ctx context.Context, dg *dgo.Dgraph,
	node interface{}, uidMap map[string]UID, m *sync.Mutex) (bool, error) {

	n, err := c.mutate(ctx, dg, node.(*DupleNode), uidMap, m)
	if err != nil {
		return n, fmt.Errorf("mutateSingleDupleNode: %w", err)
	}

	return n, nil
}

// mutateStringMap loops through the given map to create a DupleNode struct.
// Once the DupleNode is made then it is passed to mutate to be upserted.
func (c *Client) mutateStringMap(ctx context.Context, dg *dgo.Dgraph,
	d map[string]string, uidMap map[string]UID, m sync.Locker) (bool, error) {
	// Convert out map[string]string to usable predicate and value data.
	predVals := c.mapToPredValPairs(d)

	n, err := c.mutate(ctx, dg, predVals, uidMap, m)
	if err != nil {
		return n, fmt.Errorf("mutateStringMap: %w", err)
	}

	return n, nil
}

// mutateDynamicMap loops through the given map to create a DupleNode struct.
// Once the DupleNode is made then it is passed to mutate to be upserted.
// Note: The looping process supports multiple Dgraph datatypes.
func (c *Client) mutateDynamicMap(ctx context.Context, dg *dgo.Dgraph,
	d map[string]interface{}, uidMap map[string]UID, m sync.Locker) (bool, error) {
	// Convert out map[string]interface{} to usable predicate and value data.
	predVals := c.dynamicMapToPredValPairs(d)

	n, err := c.mutate(ctx, dg, predVals, uidMap, m)
	if err != nil {
		return n, fmt.Errorf("mutateDynamicMap: %w", err)
	}

	return n, nil
}
