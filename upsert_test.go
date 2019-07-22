package quirk

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"

	. "github.com/franela/goblin"
)

func TestTryUpsert(t *testing.T) {
	g := Goblin(t)
	c := NewClient()

	g.Describe("tryUpsert", func() {
		var dgraph *dgo.Dgraph
		ctx := context.Background()

		g.It("Should not error", func() {
			dgraph = dgo.NewDgraphClient(&testDgraphClient{
				queryResponse: testValidJSONOutput,
			})
			g.Assert(c.tryUpsert(ctx, dgraph.NewTxn(), testPredValCorrect)).
				Equal(&upsertResponse{identifier: blankDefault, uid: "0x1"})
		})

		g.It("Should error from mutation", func() {
			dgraph = dgo.NewDgraphClient(&testDgraphClient{
				failQueryOn:   2,
				queryResponse: []byte(`{}`),
			})
			res := c.tryUpsert(ctx, dgraph.NewTxn(), testPredValCorrect)

			if res.err == nil {
				g.Fail(res.err)
			}
		})
	})
}
