package quirk

import (
	"context"
	"sync"

	"github.com/dgraph-io/dgo"
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

type mutateSingle func(context.Context, *dgo.Dgraph, interface{}, map[string]string, *sync.Mutex) error

type rdfMode int

type PredValDat struct {
	Predicate string
	Value     interface{}
}
