package quirk

import (
	"context"
	"sync"
	"testing"

	"github.com/dgraph-io/dgo"
	. "github.com/franela/goblin"
)

func TestMutateSingleDupleNode(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		duple := &DupleNode{
			Identifier: "damienstamates",
			Duples:     []Duple{},
		}

		uidMap := make(map[string]string)
		ctx := context.Background()

		g.It("should be empty, false, and nil", func() {
			new, err := c.mutateSingleDupleNode(ctx,
				dgo.NewDgraphClient(&testDgraphClient{
					queryResponse: testValidJSONOutput,
				}),
				duple, uidMap, &sync.Mutex{})

			g.Assert(new).
				Equal(true)

			g.Assert(err).
				Equal(error(nil))

			g.Assert(len(uidMap)).
				Equal(1)
		})
	})
}
