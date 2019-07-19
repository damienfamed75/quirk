package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateSingleDupleNode(ctx context.Context, dg DgraphClient,
	node *DupleNode, uidMap map[string]string, m sync.Locker) (bool, error) {

	// predVals := reflectDupleMaps(node)

	return false, nil
}

func reflectDupleMaps(node *DupleNode) predValPairs {
	// Maybe replace all instances of predValDat+predValPairs with DupleNode
	// Maybe remove Name string because it's difficult to make?
	// Or maybe keep it for the user?

	return nil
}
