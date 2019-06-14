package quirk

import (
	"context"
	"errors"
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

	var builder strings.Builder

	uid, err := queryUID(ctx, txn, &builder, dat)
	if err != nil {
		return &upsertResponse{err: err}
	}

	identifier := blankDefault
	for _, d := range dat {
		if d.predicate == c.predicateKey {
			identifier = fmt.Sprintf("%v", d.value)
		}
	}

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

func mutateNewNode(ctx context.Context, txn dgraphTxn, b builder, identifier string, dat []*predValDat) (string, error) {
	for _, d := range dat {
		fmt.Fprintf(b, rdfBase, identifier, d.predicate, d.value)
	}

	mu := &api.Mutation{SetNquads: []byte(b.String())}
	assigned, err := txn.Mutate(ctx, mu)
	if err != nil {
		return "", err
	}

	uid := assigned.GetUids()[identifier]
	if uid == "" {
		return "", errors.New("UID not received")
	}

	return uid, nil
}
