package quirk

import (
	"context"
	"io"
	"sync"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

// Exported structures for the Client to use.
type (
	// Operation is the main parameter used when calling quirk client methods.
	Operation struct {
		SetMultiStruct     []interface{}
		SetSingleStruct    interface{}
		SetStringMap       map[string]string
		SetDynamicMap      map[string]interface{}
		SetSingleDupleNode *DupleNode
		SetMultiDupleNode  []DupleNode
	}

	DupleNode struct {
		Name   string
		Duples []Duple
	}

	// Duple is a structural way of giving the quirk client enough information
	// about a node to create triples and insert them into Dgraph.
	Duple struct {
		// Predicate acts as a key.
		Predicate string
		// Object is the data representing the predicate.
		Object interface{}
		// IsUnique stores whether or not to treat this as an upsert or not.
		IsUnique bool
	}

	// DgraphClient is used to mock out the client when testing.
	DgraphClient interface {
		Alter(context.Context, *api.Operation) error
		NewTxn() *dgo.Txn
	}
)

// non exported structures.
type (
	predValDat struct {
		predicate string
		value     interface{}
		isUnique  bool
	}
	upsertResponse struct {
		new        bool
		err        error
		identifier string
		uid        string
	}
)

// interfaces used within for testing.
type (
	dgraphTxn interface {
		Query(context.Context, string) (*api.Response, error)
		Mutate(context.Context, *api.Mutation) (*api.Assigned, error)
		Commit(context.Context) error
		Discard(context.Context) error
	}
	builder interface {
		io.Writer
		String() string
		Reset()
	}
)

// predValPairs is used to sort out the upsert valued
// predValDat from the slice.
type predValPairs []*predValDat

func (p predValPairs) unique() (pairs predValPairs) {
	for _, v := range p {
		if v.isUnique {
			pairs = append(pairs, v)
		}
	}
	return pairs
}

// Credit: The Go Authors @ "encoding/json"
// tagOptions is the string following a comma in a struct field's "quirk"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// queryDecode is our type when unmarshalling a query response.
type queryDecode map[string][]struct{ UID *string }

// mutateSingle is used to pass into a worker function to call.
type mutateSingle func(context.Context, DgraphClient, interface{}, map[string]string, *sync.Mutex) (bool, error)
