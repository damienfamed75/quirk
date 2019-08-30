package quirk

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"

	. "github.com/franela/goblin"
)

func TestSetNode(t *testing.T) {
	g := Goblin(t)

	g.Describe("setNode", func() {
		dgraph := dgo.NewDgraphClient(&testDgraphClient{
			queryResponse: testValidJSONOutput,
		})
		ctx := context.Background()

		g.It("Should not error", func() {
			s, err := setNode(ctx, dgraph.NewTxn(),
				&testBuilder{}, "damienstamates", testPredValCorrect)

			g.Assert(s["damienstamates"]).Equal("0x1")
			g.Assert(err).Equal(error(nil))
		})
	})
}
