package quirk

import (
	"context"
	"sync"

	"github.com/damienfamed75/quirk/logging"
)

// Client is used to store enough data and help manage
// the logger when inserting nodes into Dgraph using a proper
// upsert procedure.
type Client struct {
	predicateKey string
	logger       logging.Logger
}

// setupClient returns the default states of a quirk client.
func setupClient() *Client {
	return &Client{
		logger:       NewNilLogger(),
		predicateKey: "name",
	}
}

// NewClient will setup a new client with the passed in
// configurations if so chosen to use any.
func NewClient(confs ...ClientConfiguration) *Client {
	q := setupClient()

	// Loop through the configurations and apply them to the client.
	for _, c := range confs {
		c(q)
	}

	return q
}

// InsertNode takes in an Operation to determine if multiple nodes
// will be added or a single node. Then the function will return a
// map of the returned successful UIDs with the key being the predicate
// key value. By default this will be the "name" predicate value.
func (c *Client) InsertNode(ctx context.Context, dg DgraphClient, o *Operation) (uidMap map[string]string, err error) {
	if o.SetMultiStruct != nil && o.SetSingleStruct != nil {
		return nil, &Error{
			Msg:      msgTooManyMutationFields,
			File:     "client.go",
			Function: "quirk.Client.InsertNode",
		}
	}

	uidMap = make(map[string]string)

	if o.SetMultiStruct != nil {
		err = c.mutateMultiStruct(ctx, dg, o.SetMultiStruct, uidMap)
	} else if o.SetSingleStruct != nil {
		_, err = c.mutateSingleStruct(ctx, dg, o.SetSingleStruct, uidMap, &sync.Mutex{})
	} else if o.SetStringMap != nil {
		_, err = c.mutateStringMap(ctx, dg, o.SetStringMap, uidMap, &sync.Mutex{})
	} else if o.SetDynamicMap != nil {
		_, err = c.mutateDynamicMap(ctx, dg, o.SetDynamicMap, uidMap, &sync.Mutex{})
	}

	return
}

// GetPredicateKey returns the name of the field(predicate) that will
// be used to label inserted nodes. By default this is "name"
func (c *Client) GetPredicateKey() string {
	return c.predicateKey
}
