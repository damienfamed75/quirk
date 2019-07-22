package quirk

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"

	. "github.com/franela/goblin"
)

func TestSetNewNode(t *testing.T) {
	g := Goblin(t)

	g.Describe("tryUpsert", func() {
		dgraph := dgo.NewDgraphClient(&testDgraphClient{
			queryResponse: testValidJSONOutput,
		})
		ctx := context.Background()

		g.It("Should not error", func() {
			s, err := setNewNode(ctx, dgraph.NewTxn(),
				&testBuilder{}, "damienstamates", testPredValCorrect)

			g.Assert(s).Equal("0x1")
			g.Assert(err).Equal(error(nil))
		})

		g.It("Should error", func() {
			s, err := setNewNode(ctx, dgraph.NewTxn(),
				&testBuilder{}, "fail", testPredValCorrect)

			g.Assert(s).Equal("")
			g.Assert(err).Equal(&TransactionError{
				Msg: msgMutationHadNoUID, Function: "setNewNode"})
		})
	})
}
