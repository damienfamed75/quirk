package quirk

import (
	"context"
	"testing"

	"github.com/dgraph-io/dgo"
	. "github.com/franela/goblin"
	"go.uber.org/zap/zapcore"
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
	c := NewClient(WithPredicateKey("username"), WithLogger(NewCustomLogger([]byte("dpanic"), zapcore.EncoderConfig{})))
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

	g.Describe("String Map", func() {

		stringMap := make(map[string]string)
		stringMap["username"] = "damienstamates"
		stringMap["website"] = "github.com"
		stringMap["accountAge"] = "197"
		stringMap["email"] = "damienstamates@gmail.com"

		g.It("should not error and return an empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{}),
				&Operation{SetStringMap: stringMap})

			g.Assert(len(uids)).
				Equal(1)

			g.Assert(err).
				Equal(error(nil))
		})
	})

	g.Describe("Dynamic Map", func() {

		dynamicMap := make(map[string]interface{})
		dynamicMap["username"] = "damienstamates"
		dynamicMap["website"] = "github.com"
		dynamicMap["accountAge"] = 197
		dynamicMap["email"] = "damienstamates@gmail.com"

		g.It("should not error and return an empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{}),
				&Operation{SetDynamicMap: dynamicMap})

			g.Assert(len(uids)).
				Equal(1)

			g.Assert(err).
				Equal(error(nil))
		})
	})

	g.Describe("Duple Node", func() {

		duple := &DupleNode{
			Identifier: "damienstamates",
			Duples: []Duple{
				Duple{Predicate: "username", Object: "damienstamates"},
				Duple{Predicate: "website", Object: "github.com"},
				Duple{Predicate: "accountAge", Object: 197},
				Duple{Predicate: "email", Object: "damienstamates@gmail.com"},
			},
		}

		g.It("should not error and return an empty map", func() {
			uids, err := c.InsertNode(ctx, dgo.NewDgraphClient(&testDgraphClient{}),
				&Operation{SetSingleDupleNode: duple})

			g.Assert(len(uids)).
				Equal(1)

			g.Assert(err).
				Equal(error(nil))
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
