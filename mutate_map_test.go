package quirk

import (
	"context"
	"sync"
	"testing"

	"github.com/dgraph-io/dgo"
	. "github.com/franela/goblin"
)

func TestStringMap(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		testMap := make(map[string]string)
		testMap["username"] = "damienstamates"
		testMap["website"] = "github.com"
		testMap["accountAge"] = "197"
		testMap["email"] = "damienstamates@gmail.com"

		uidMap := make(map[string]string)
		ctx := context.Background()

		g.It("should be empty, false, and nil", func() {
			new, err := c.mutateStringMap(ctx,
				dgo.NewDgraphClient(&testDgraphClient{
					queryResponse: testValidJSONOutput,
				}),
				testMap, uidMap, &sync.Mutex{})

			g.Assert(new).
				Equal(true)

			g.Assert(err).
				Equal(error(nil))

			g.Assert(len(uidMap)).
				Equal(1)
		})
	})
}

func TestDynamicMap(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		testMap := make(map[string]interface{})
		testMap["username"] = "damienstamates"
		testMap["website"] = "github.com"
		testMap["accountAge"] = 197
		testMap["email"] = "damienstamates@gmail.com"

		uidMap := make(map[string]string)
		ctx := context.Background()

		g.It("should be empty, false, and nil", func() {
			new, err := c.mutateDynamicMap(ctx,
				dgo.NewDgraphClient(&testDgraphClient{
					queryResponse: testValidJSONOutput,
				}),
				testMap, uidMap, &sync.Mutex{})

			g.Assert(new).
				Equal(true)

			g.Assert(err).
				Equal(error(nil))

			g.Assert(len(uidMap)).
				Equal(1)
		})
	})
}
