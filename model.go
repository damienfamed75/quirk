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
		SetMultiStruct  []interface{}
		SetSingleStruct interface{}
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
		Predicate string
		Value     interface{}
		IsUpsert  bool
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

// tagOptions is used to identify the type of an optional quirk tag.
type tagOptions string

// queryDecode is our type when unmarshalling a query response.
type queryDecode map[string][]struct{ UID *string }

// mutateSingle is used to pass into a worker function to call.
type mutateSingle func(context.Context, DgraphClient, interface{}, map[string]string, *sync.Mutex) (bool, error)