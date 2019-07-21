package quirk

// import (
// 	"context"
// 	"sync"
// 	"testing"

// 	"github.com/dgraph-io/dgo"
// 	. "github.com/franela/goblin"
// )

// func TestMutateSingleStruct(t *testing.T) {
// 	g := Goblin(t)
// 	c := NewClient(WithPredicateKey("username"))

// 	g.Describe("UIDMap, New, and Error", func() {
// 		uidMap := make(map[string]string)
// 		ctx := context.Background()

// 		g.It("should be empty, false, and nil", func() {
// 			new, err := c.mutateSingleStruct(ctx,
// 				dgo.NewDgraphClient(&testDgraphClient{queryResponse: testValidJSONOutput}),
// 				&testPersonCorrect, uidMap, &sync.Mutex{})

// 			g.Assert(new).
// 				Equal(false)

// 			g.Assert(err).
// 				Equal(error(nil))

// 			g.Assert(len(uidMap)).
// 				Equal(0)
// 		})

// 		g.It("should be not empty, false, and nil when new", func() {
// 			new, err := c.mutateSingleStruct(ctx,
// 				dgo.NewDgraphClient(&testDgraphClient{
// 					queryResponse: []byte("{}"),
// 				}),
// 				&testPersonCorrect, uidMap, &sync.Mutex{})

// 			g.Assert(new).
// 				Equal(true)

// 			g.Assert(err).
// 				Equal(error(nil))

// 			g.Assert(len(uidMap)).
// 				Equal(1)

// 			g.Assert(uidMap["damienstamates"]).
// 				Equal("0x1")
// 		})
// 	})
// }

// func TestTryUpsert(t *testing.T) {
// 	g := Goblin(t)
// 	c := NewClient()

// 	g.Describe("tryUpsert", func() {
// 		ctx := context.Background()

// 		g.It("Should not error", func() {
// 			g.Assert(c.tryUpsert(ctx, &testTxn{jsonOutput: testValidJSONOutput}, testPredValCorrect)).
// 				Equal(&upsertResponse{identifier: blankDefault, uid: "0x1"})
// 		})

// 		g.It("Should error from mutation", func() {
// 			res := c.tryUpsert(ctx, &testTxn{
// 				failOn: 2, jsonOutput: []byte("{}")}, testPredValCorrect)

// 			if res.err == nil {
// 				g.Fail(res.err)
// 			}
// 		})
// 	})
// }

// func TestSetNewNode(t *testing.T) {
// 	g := Goblin(t)

// 	g.Describe("tryUpsert", func() {
// 		ctx := context.Background()

// 		g.It("Should not error", func() {
// 			s, err := setNewNode(ctx, &testTxn{
// 				jsonOutput: testValidJSONOutput},
// 				&testBuilder{}, "a", testPredValCorrect)

// 			g.Assert(s).Equal("0x1")
// 			g.Assert(err).Equal(error(nil))
// 		})

// 		g.It("Should error", func() {
// 			s, err := setNewNode(ctx, &testTxn{
// 				jsonOutput: testValidJSONOutput},
// 				&testBuilder{}, "b", testPredValCorrect)

// 			g.Assert(s).Equal("")
// 			g.Assert(err).Equal(&TransactionError{
// 				Msg: msgMutationHadNoUID, Function: "mutateNewNode"})
// 		})
// 	})
// }
