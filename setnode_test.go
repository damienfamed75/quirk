package quirk

import (
	"context"
	"strings"
	"testing"

	"github.com/dgraph-io/dgo/v2"

	. "github.com/franela/goblin"
)

func TestSetNode(t *testing.T) {
	g := Goblin(t)

	g.Describe("tryUpsert", func() {
		dgraph := dgo.NewDgraphClient(&testDgraphClient{
			queryResponse: testValidJSONOutput,
		})
		ctx := context.Background()

		g.It("Should not error", func() {
			var builder strings.Builder
			s, err := setNode(ctx, dgraph.NewTxn(),
				&builder, "_:damienstamates", testPredValCorrect)

			g.Assert(s["damienstamates"]).Equal("0x1")
			g.Assert(err).Equal(error(nil))
		})
	})
}
