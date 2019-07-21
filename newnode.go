package quirk

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/protos/api"
)

// setNewNode will build a mutation RDF with the builder and will then
// execute it using the given transaction. Once executed it will return the UID.
func setNewNode(ctx context.Context, txn dgraphTxn, b builder,
	identifier string, dat *DupleNode) (string, error) {
	for _, d := range dat.Duples {
		// get any optional XML datatype knowledge based on the value.
		fmt.Fprintf(b, rdfBase+d.dataType+rdfEnd, identifier, d.Predicate, d.Object)
	}

	// Use our transaction to execute a mutation to add our new node.
	assigned, err := txn.Mutate(ctx, &api.Mutation{SetNquads: []byte(b.String())})
	if err != nil {
		return "", err
	}

	uid := assigned.GetUids()[identifier]
	if uid == "" {
		return "", &TransactionError{
			Msg: msgMutationHadNoUID, Function: "setNewNode", RDF: b.String()}
	}

	return uid, nil
}
