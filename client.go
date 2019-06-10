package quirk

import (
	"context"
	"sync"

	"github.com/damienfamed75/quirk/logging"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"go.uber.org/zap"
)

type Client struct {
	schemaCache Schema

	quirkName      string
	quirkID        uint64
	quirkRel       string
	schemaString   string
	useIncrementor bool
	predicateKey   string
	quirkReverse   bool

	logger logging.Logger
}

func setupClient() *Client {
	return &Client{
		schemaCache:    make(map[string]Properties),
		useIncrementor: false,
		logger:         NewNilLogger(),
		predicateKey:   "name",
		quirkName:      "quirk",
		quirkRel:       "hasExact",
		quirkReverse:   false,
	}
}

func NewClient(schema string, confs ...ClientConfiguration) (*Client, error) {
	q := setupClient()

	for _, c := range confs {
		c(q)
	}

	err := q.setSchema(schema)
	q.quirkID = aeshash(q.quirkName)
	if err != nil {
		q.logger.Warn("Schema was not processed correctly.", zap.Error(err))
	}

	return q, err
}

func (c *Client) InitializeSchema(ctx context.Context, dg *dgo.Dgraph) error {
	return dg.Alter(ctx, &api.Operation{Schema: c.schemaString})
}

func (c *Client) CreateRDF(m *Options) string {
	// Check if both values are not nil

	// Create RDF the given struct/structs

	return ""
}

func (c *Client) InsertNode(ctx context.Context, dg *dgo.Dgraph, o *Options) (uidMap map[string]string, err error) {
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
		err = c.mutateSingleStruct(ctx, dg, o.SetSingleStruct, uidMap, &sync.Mutex{})
	}

	return
}

func (c *Client) GetSchema() Schema {
	return c.schemaCache
}

func (c *Client) GetPredicateKey() string {
	return c.predicateKey
}
