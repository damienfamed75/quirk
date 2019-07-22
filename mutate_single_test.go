package quirk

import (
	"context"
	"sync"
	"testing"

	"github.com/dgraph-io/dgo"
	. "github.com/franela/goblin"
)

func TestMutateSingleStruct(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		uidMap := make(map[string]string)
		ctx := context.Background()

		g.It("should be empty, false, and nil", func() {
			new, err := c.mutateSingleStruct(ctx,
				dgo.NewDgraphClient(&testDgraphClient{queryResponse: testValidJSONOutput}),
				&testPersonCorrect, uidMap, &sync.Mutex{})

			g.Assert(new).
				Equal(false)

			g.Assert(err).
				Equal(error(nil))

			g.Assert(len(uidMap)).
				Equal(0)
		})

		g.It("should be not empty, false, and nil when new", func() {
			new, err := c.mutateSingleStruct(ctx,
				dgo.NewDgraphClient(&testDgraphClient{
					queryResponse: []byte("{}"),
				}),
				&testPersonCorrect, uidMap, &sync.Mutex{})

			g.Assert(new).
				Equal(true)

			g.Assert(err).
				Equal(error(nil))

			g.Assert(len(uidMap)).
				Equal(1)

			g.Assert(uidMap["damienstamates"]).
				Equal("0x1")
		})
	})
}
