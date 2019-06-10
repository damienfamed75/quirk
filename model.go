package quirk

import (
	"context"
	"sync"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

type Options struct {
	SetMultiStruct  []interface{}
	SetSingleStruct interface{}
}

type Schema map[string]Properties

type Properties struct {
	Upsert  bool
	Lang    bool
	Reverse bool
	Index   indexType
	DType   string
}

type indexType struct {
	IsIndex bool
	DType   string
}

type mutateSingle func(context.Context, DgraphClient, interface{}, map[string]string, *sync.Mutex) error

type rdfMode int

type PredValDat struct {
	Predicate string
	Value     interface{}
}

type DgraphClient interface {
	Alter(context.Context, *api.Operation) error
	NewTxn() *dgo.Txn
}

type DgraphTxn interface {
	Mutate(context.Context, *api.Mutation) (*api.Assigned, error)
	Commit(context.Context) error
	Discard(context.Context) error
}
