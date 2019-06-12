package quirk

import (
	"context"
	"sync"

	"github.com/damienfamed75/quirk/logging"
)

type Client struct {
	predicateKey string
	logger logging.Logger
}

func setupClient() *Client {
	return &Client{
		logger:       NewNilLogger(),
		predicateKey: "name",
	}
}

func NewClient(confs ...ClientConfiguration) *Client {
	q := setupClient()

	for _, c := range confs {
		c(q)
	}

	return q
}

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
	}

	return
}

func (c *Client) GetPredicateKey() string {
	return c.predicateKey
}
