package quirk

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

func (c *Client) mutateSingleStruct(ctx context.Context, dg *dgo.Dgraph, d interface{}, uidMap map[string]string, m *sync.Mutex) error {
	upsertMap, fullMap := c.reflectMaps(d)

	// send upsert separate map to get hashed UIDs+RDF.
	upsertRDFs := c.createRDF(hash, upsertMap)

	fmt.Printf("\n\nUPPPP: [\n%s]", upsertRDFs)
	// send upsert RDF to be put into dgraph.
	// Note: We're not saving the upsert uid map, because the user
	// doesn't need to know these values. They just need to know the
	// values of the nodes they asked to insert.
	// failedUIDMap := make(map[string]string)
	for _, rdf := range upsertRDFs {
		err := mutate(ctx, dg.NewTxn(), rdf, make(map[string]string), m)
		if err != nil {
			return &FailedUpsert{PredMap: upsertMap, RDF: rdf}
		}
	}

	// Get out designated mode for how the UIDs are going to be
	// generated with our actual data nodes.
	var mode = auto
	if c.useIncrementor {
		mode = incrementor
	}

	// send full map to get auto incremented RDF string.
	fullRDFs := c.createRDF(mode, fullMap)

	// else mutate the second RDF of nonupsert.
	err := mutate(ctx, dg.NewTxn(), strings.Join(fullRDFs, "\n"), uidMap, m)
	if err != nil {
		return err
	}
	fmt.Printf("\nFULLL: [\n%s]\n\n", fullRDFs)

	// return UIDs of second RDF's nodes.
	return nil
}

func (c *Client) reflectMaps(d interface{}) (upsert map[string]interface{}, full map[string]interface{}) {
	upsert = make(map[string]interface{})
	full = make(map[string]interface{})
	var elem = reflect.ValueOf(d).Elem()

	// loop through elements of struct.
	for i := 0; i < elem.NumField(); i++ {
		var tag = reflect.TypeOf(d).Elem().Field(i).Tag.Get("quirk")
		// store upsert predicates in separate map.
		if c.isUpsert(tag) {
			// If this is an upsert then add it to the upsert
			// to be treated specially.
			upsert[tag] = elem.Field(i).Interface()
		}
		// Add the predicate and value to the full map.
		full[tag] = elem.Field(i).Interface()
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
