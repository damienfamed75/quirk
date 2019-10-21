package quirk

import (
	"context"
	"fmt"
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
		if qerr, ok := err.(*QueryError); ok {
			return &upsertResponse{err: qerr.setVerbose(c.verbose)}
		}
		return &upsertResponse{err: err}
	}

	// Check if the given data contains the quirk client predicateKey. If not
	// then it is defaulted to "data"
	identifier := _blankDefault
	if dat.Identifier != "" {
		identifier = dat.Identifier
	}

	// If the UID was not found by our query then mutate to add a new node.
	var new bool
	if uid == "" {
		new = true
		// Insert new node.
		uidMap, err := setNode(ctx, txn, &builder, "_:"+identifier, dat)
		if err != nil {
			if merr, ok := err.(*MutationError); ok {
				return &upsertResponse{
					err: merr.setVerbose(c.verbose).setNew(new), new: new}
			}
			// Even though this function can't return any other types of errors
			// if in the future I add a new error type, then this will help
			// with that.
			return &upsertResponse{
				err: fmt.Errorf("setNode (insert): %w", err),
				new: new,
			}
		}
		// If the UID could not be found in the map.
		if uid = uidMap[identifier]; uid == "" {
			return &upsertResponse{
				err: ErrUIDNotFound,
				new: new,
			}
		}
	} else {
		// Update the found node.
		_, err = setNode(ctx, txn, &builder, "<"+uid+">", dat)
		if err != nil {
			if merr, ok := err.(*MutationError); ok {
				return &upsertResponse{
					err: merr.setVerbose(c.verbose).setNew(new), new: new}
			}
			// Even though this function can't return any other types of errors
			// if in the future I add a new error type, then this will help
			// with that.
			return &upsertResponse{
				err: fmt.Errorf("setNode (update): %w", err),
				new: new,
			}
		}
	}

	return &upsertResponse{
		err: txn.Commit(ctx), new: new, uid: uid, identifier: identifier,
	}
}
