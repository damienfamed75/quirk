package quirk

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

// setNode will build a mutation RDF with the builder and will then
// execute it using the given transaction. Once executed it will return the UID map.
func setNode(ctx context.Context, txn *dgo.Txn, b builder,
	identifier string, dat *DupleNode) (map[string]string, error) {
	for _, d := range dat.Duples {
		d.dataType = checkType(d.Object)
		if uid, ok := d.Object.(UID); ok {
			// Use the UID format instead of the regular object.
			fmt.Fprintf(b, rdfReference+d.dataType+rdfEnd, identifier, d.Predicate, uid.Value())
		} else if slice, ok := d.Object.([]byte); ok {
			// get any optional XML datatype knowledge based on the value.
			fmt.Fprintf(b, rdfBase+rdfEnd, identifier, d.Predicate, string(slice))
		} else {
			// get any optional XML datatype knowledge based on the value.
			fmt.Fprintf(b, rdfBase+d.dataType+rdfEnd, identifier, d.Predicate, d.Object)
		}
	}

	// Use our transaction to execute a mutation to add our new node.
	assigned, err := txn.Mutate(ctx, &api.Mutation{SetNquads: []byte(b.String())})
	if err != nil {
		return nil, err
	}

	return assigned.GetUids(), nil
}
