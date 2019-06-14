package quirk

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"
	. "github.com/franela/goblin"
)

func BenchmarkSetupClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = setupClient()
	}
}

func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewClient(WithLogger(NewDebugLogger()))
	}
}

func BenchmarkInsertNode(b *testing.B) {
	c := NewClient()
	ctx := context.Background()
	dg := dgo.NewDgraphClient(&testDgraphClient{
		queryResponse: testValidJSONOutput})
	op := &Operation{SetSingleStruct: &testPersonCorrect}

	for i := 0; i < b.N; i++ {
		_, _ = c.InsertNode(ctx, dg, op)
	}
}

func BenchmarkGetPredicateKey(b *testing.B) {
	c := NewClient()
	for i := 0; i < b.N; i++ {
		_ = c.GetPredicateKey()
	}
}

func TestSetupClient(t *testing.T) {
	g := Goblin(t)

	g.Describe("Client", func() {
		g.It("should have nil logger", func() {
			g.Assert(setupClient().logger).
				Equal(NewNilLogger())
		})
		g.It("predicate key should be \"name\"", func() {
			g.Assert(setupClient().predicateKey).
				Equal("name")
		})
	})
}

func TestNewClient(t *testing.T) {
	g := Goblin(t)

	g.Describe("NewClient", func() {
		g.It("should have debug logger", func() {
			dl := NewDebugLogger()
			g.Assert(NewClient(WithLogger(dl)).logger).
				Equal(dl)
		})
	})
}

func TestInsertNode(t *testing.T) {
	g := Goblin(t)
	c := NewClient()
	ctx := context.Background()

	g.Describe("Single Struct", func() {

		g.It("should not error and return an empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{
				queryResponse: testValidJSONOutput}),
				&Operation{SetSingleStruct: &testPersonCorrect})

			g.Assert(len(uids)).
				Equal(0)

			g.Assert(err).
				Equal(error(nil))
		})

		g.It("should error and return empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{}),
				&Operation{SetSingleStruct: &testPersonCorrect})

			g.Assert(len(uids)).
				Equal(0)

			if err == nil {
				g.Fail(err)
			}
		})

	})

	g.Describe("Multiple Structs", func() {

		multiset := []interface{}{&testPersonCorrect}

		g.It("should not error and return an empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{
				queryResponse: testValidJSONOutput}),
				&Operation{SetMultiStruct: multiset})

			g.Assert(len(uids)).
				Equal(0)

			g.Assert(err).
				Equal(error(nil))
		})

		g.It("should error and return empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{}),
				&Operation{SetMultiStruct: multiset})

			g.Assert(len(uids)).
				Equal(0)

			if err == nil {
				g.Fail(err)
			}
		})

	})

	g.Describe("Single and multi structs", func() {

		g.It("should return an error", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{}),
				&Operation{SetMultiStruct: []interface{}{}, SetSingleStruct: 1})

			g.Assert(len(uids)).
				Equal(0)

			g.Assert(err).
				Equal(&Error{
					Msg:      msgTooManyMutationFields,
					File:     "client.go",
					Function: "quirk.Client.InsertNode",
				})
		})
	})
}

func TestGetPredicateKey(t *testing.T) {
	g := Goblin(t)

	g.Describe("Predicate key", func() {
		c := NewClient()
		g.It("should equal, \"c.predicateKey\"", func() {
			g.Assert(c.GetPredicateKey()).
				Equal(c.predicateKey)
		})
	})

}
