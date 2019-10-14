package quirk

import (
	"context"
	"sync"
	"testing"

	"github.com/dgraph-io/dgo/v2"
	. "github.com/franela/goblin"
)

func TestMutateSingleStruct(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		uidMap := make(map[string]UID)
		ctx := context.Background()

		g.It("should be empty, false, and nil", func() {
			new, err := c.mutateSingleStruct(ctx,
				dgo.NewDgraphClient(&testDgraphClient{queryResponse: testValidJSONOutput}),
				&testPersonCorrect, uidMap, &sync.Mutex{})

			g.Assert(new).
				Equal(false)

			g.Assert(err).
				Equal(error(nil))
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

			g.Assert(uidMap["damienstamates"].uid).
				Equal("0x1")
		})
	})
}

func TestMutateSingleDupleNode(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		duple := &DupleNode{
			Identifier: "damienstamates",
			Duples:     []Duple{},
		}

		uidMap := make(map[string]UID)
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

func TestStringMap(t *testing.T) {
	g := Goblin(t)
	c := NewClient(WithPredicateKey("username"))

	g.Describe("UIDMap, New, and Error", func() {
		testMap := make(map[string]string)
		testMap["username"] = "damienstamates"
		testMap["website"] = "github.com"
		testMap["accountAge"] = "197"
		testMap["email"] = "damienstamates@gmail.com"

		uidMap := make(map[string]UID)
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

		uidMap := make(map[string]UID)
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
