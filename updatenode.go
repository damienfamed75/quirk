package quirk

import (
	"context"
	"fmt"

	"github.com/dgraph-io/dgo/protos/api"
)

func updateNode(ctx context.Context, txn dgraphTxn, b builder,
	identifier string, dat *DupleNode, oldUID string) error {
	for _, d := range dat.Duples {
		d.dataType = checkType(d.Object)
		if uid, ok := d.Object.(UID); ok {
			// Use the UID format instead of the regular object.
			fmt.Fprintf(b, rdfOldReference+d.dataType+rdfEnd, oldUID, d.Predicate, uid.uid)
		} else if slice, ok := d.Object.([]byte); ok {
			// get any optional XML datatype knowledge based on the value.
			fmt.Fprintf(b, rdfOldBase+rdfEnd, oldUID, d.Predicate, string(slice))
		} else {
			// get any optional XML datatype knowledge based on the value.
			fmt.Fprintf(b, rdfOldBase+d.dataType+rdfEnd, oldUID, d.Predicate, d.Object)
		}
	}

	_, err := txn.Mutate(ctx, &api.Mutation{SetNquads: []byte(b.String())})
	if err != nil {
		return err
	}

	return nil
}
