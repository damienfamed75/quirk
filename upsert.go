package quirk

import (
	"context"
	"errors"
	"strings"

	"github.com/dgraph-io/dgo/v2"
)

func (c *Client) tryUpsert(ctx context.Context, txn *dgo.Txn, dat *DupleNode) *upsertResponse {
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
	if dat.Identifier != "" {
		identifier = dat.Identifier
	}

	// If the UID was not found by our query then mutate to add a new node.
	var new bool
	if uid == "" {
		new = true
		// Insert new node.
		// TODO remove string concatenation.
		uidMap, err := setNode(ctx, txn, &builder, "_:"+identifier, dat)
		if err != nil {
			return &upsertResponse{
				err: err,
				new: new,
			}
		}
		// If the UID could not be found in the map.
		if uid = uidMap[identifier]; uid == "" {
			return &upsertResponse{
				err: errors.New(msgMutationHadNoUID),
				new: new,
			}
		}
	} else {
		// Update the found node.
		// TODO remove string concatenation.
		_, err = setNode(ctx, txn, &builder, "<"+uid+">", dat)
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
