package quirk

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/dgraph-io/dgo/protos/api"
)

func (c *Client) mutateSingleStruct(ctx context.Context, dg DgraphClient,
	d interface{}, uidMap map[string]string, m *sync.Mutex) (bool, error) {
	// Use reflect to package the predicate and values in slices.
	predVals := reflectMaps(d)

	res := c.tryUpsert(ctx, dg.NewTxn(), predVals)

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
		uid, err = mutateNewNode(ctx, txn, &builder, identifier, dat)
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

// mutateNewNode will build a mutation RDF with the builder and will then
// execute it using the given transaction. Once executed it will return the UID.
func mutateNewNode(ctx context.Context, txn dgraphTxn, b builder,
	identifier string, dat []*predValDat) (string, error) {
	for _, d := range dat {
		fmt.Fprintf(b, rdfBase, identifier, d.predicate, d.value)
	}

	// Use our transaction to execute a mutation to add our new node.
	assigned, err := txn.Mutate(ctx, &api.Mutation{SetNquads: []byte(b.String())})
	if err != nil {
		return "", err
	}

	uid := assigned.GetUids()[identifier]
	if uid == "" {
		return "", &TransactionError{
			Msg: msgMutationHadNoUID, Function: "mutateNewNode", RDF: b.String()}
	}

	return uid, nil
}
