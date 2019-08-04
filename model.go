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
	// Note: only one of these should be filled at a time, because only one
	// will be executed and taken care of as seen in client.go
	Operation struct {
		SetMultiStruct     []interface{}
		SetSingleStruct    interface{}
		SetStringMap       map[string]string
		SetDynamicMap      map[string]interface{}
		SetSingleDupleNode *DupleNode
		SetMultiDupleNode  []*DupleNode
	}

	// DupleNode is the container for a duple node.
	DupleNode struct {
		Identifier string
		Duples     []Duple
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
		// dataType stores the xml tag for the datatype.
		dataType string
	}

	// DgraphClient is used to mock out the client when testing.
	DgraphClient interface {
		Alter(context.Context, *api.Operation) error
		NewTxn() *dgo.Txn
	}

	// UID is used to identify the ID's given to the user and retrieved back to
	// be put as the object of a predicate.
	// This way quirk can handle the UID how they're supposed to be handled.
	// Note: Use this struct as the Object for Duples to create relationships.
	UID struct {
		uid   string
		isNew bool
	}
)

// non exported structures.
type (
	// upsertResponse is used to cleanup the large amount of info
	// that the upserting function returns.
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

// non exported aliases.
type (
	// Credit: The Go Authors @ "encoding/json"
	// tagOptions is the string following a comma in a struct field's "quirk"
	// tag, or the empty string. It does not include the leading comma.
	tagOptions string

	// queryDecode is our type when unmarshalling a query response.
	queryDecode map[string][]struct{ UID *string }

	// mutateSingle is used to pass into a worker function to call.
	mutateSingle func(context.Context, DgraphClient, interface{}, map[string]UID, *sync.Mutex) (bool, error)
)

// DupleNode ---

// Unique will loop through the Duples and return a new slice
// containing all duples that are marked as unique.
func (d *DupleNode) Unique() (duples []Duple) {
	for i := 0; i < len(d.Duples); i++ {
		if d.Duples[i].IsUnique {
			duples = append(duples, d.Duples[i])
		}
	}
	return duples
}

// Find will return a reference to a duple given that it is found
// in the slice of duples in the DupleNode.
func (d *DupleNode) Find(predicate string) *Duple {
	for i := 0; i < len(d.Duples); i++ {
		if d.Duples[i].Predicate == predicate {
			return &d.Duples[i]
		}
	}

	return nil
}

// SetOrAdd will set a pre existing duple in the DupleNode or
// if the Duple doesn't exist, then it will be added to the Node.
func (d *DupleNode) SetOrAdd(duple Duple) *DupleNode {
	if found := d.Find(duple.Predicate); found != nil {
		found.Object = duple.Object
		found.IsUnique = duple.IsUnique
		return d
	}

	return d.AddDuples(duple)
}

// AddDuples appends new duples given in the function.
// Then returns the reference to the DupleNode.
func (d *DupleNode) AddDuples(duple ...Duple) *DupleNode {
	d.Duples = append(d.Duples, duple...)
	return d
}

// UID ---

// Value returns the raw string value of the UID.
// Note: Do not use this as the Object for Duples to create relationships.
func (u UID) Value() string {
	return u.uid
}

// IsNew returns a simple boolean value indicating whether or not this node
// was a newly added node or if it was pre existing in Dgraph.
func (u UID) IsNew() bool {
	return u.isNew
}
