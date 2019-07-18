package quirk

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) tryUpsert(ctx context.Context, txn dgraphTxn, dat predValPairs) *upsertResponse {
	defer txn.Discard(ctx)

	// Pass this builder around to other functions for less mem alloc.
	var builder strings.Builder

	// Query to find if there are pre existing nodes with the unique predicates.
	uid, err := queryUID(ctx, txn, &builder, dat)
	if err != nil {
		return &upsertResponse{err: err}
	}

	// Check if the given data contains the quirk client predicateKey. If not
	// then it is defaulted to "data"
	identifier := blankDefault
	for _, d := range dat {
		if d.predicate == c.predicateKey {
			identifier = fmt.Sprintf("%v", d.value)
		}
	}

	// If the UID was not found by our query then mutate to add a new node.
	var new bool
	if uid == "" {
		new = true
		uid, err = setNewNode(ctx, txn, &builder, identifier, dat)
		if err != nil {
			return &upsertResponse{
				err: err,
				new: new,
			}
		}
	}

	return &upsertResponse{
		err: txn.Commit(ctx), new: new, uid: uid, identifier: identifier,
	}
}
