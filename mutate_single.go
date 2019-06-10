package quirk

import (
	"context"
	"reflect"
	"sync"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

func (c *Client) mutateSingleStruct(ctx context.Context, dg DgraphClient,
	d interface{}, uidMap map[string]string, m *sync.Mutex) error {
	// Use reflect to package the predicate and values in slices.
	upsertPredVals, fullPredVals := c.reflectMaps(d)

	// send upsert separate map to get hashed UIDs+RDF.
	upsertRDFs := c.createRDF(hash, upsertPredVals)

	// send upsert RDF to be put into dgraph.
	// Note: the upserted uids are not returned because the user only
	// needs the returned uids of the nodes they requested to be inserted.
	err := mutate(ctx, dg.NewTxn(), upsertRDFs, make(map[string]string), m)
	if err != nil {
		return &FailedUpsert{PredVals: upsertPredVals, RDF: upsertRDFs}
	}

	// send full map to get auto incremented RDF string.
	fullRDFs := c.createRDF(c.insertMode, fullPredVals)

	// else mutate the second RDF of nonupsert.
	err = mutate(ctx, dg.NewTxn(), fullRDFs, uidMap, m)
	if err != nil {
		return &Error{File: "mutate_single.go", Function: "mutateSingleStruct",
			Msg: msgTransactionFailure, ExtErr: err}
	}

	// return UIDs of second RDF's nodes.
	return nil
}

func (c *Client) reflectMaps(d interface{}) (upsert []*PredValDat, full []*PredValDat) {
	var elem = reflect.ValueOf(d).Elem()
	upsert = make([]*PredValDat, elem.NumField())
	full = make([]*PredValDat, elem.NumField())

	// loop through elements of struct.
	for i := 0; i < elem.NumField(); i++ {
		var tag = reflect.TypeOf(d).Elem().Field(i).Tag.Get("quirk")
		// store upsert predicates in separate map.
		if c.isUpsert(tag) {
			// If this is an upsert then add it to the upsert
			// to be treated specially.
			upsert[i] = &PredValDat{Predicate: tag, Value: elem.Field(i).Interface()}
		}
		// Add the predicate and value to the full map.
		full[i] = &PredValDat{Predicate: tag, Value: elem.Field(i).Interface()}
	}

	return
}

func (c *Client) isUpsert(tag string) bool {
	if v, ok := c.schemaCache[tag]; ok {
		return v.Upsert
	}
	return false
}

func mutate(ctx context.Context, t *dgo.Txn, rdf string, uidMap map[string]string, m *sync.Mutex) error {
	a, err := t.Mutate(ctx, &api.Mutation{
		CommitNow: true,
		SetNquads: []byte(rdf),
	})
	if err != nil {
		return err
	}

	if err = t.Discard(ctx); err != nil {
		return err
	}

	m.Lock()
	for k, v := range a.GetUids() {
		uidMap[k] = v
	}
	m.Unlock()

	return nil
}
