package quirk

import (
	"context"
	"errors"
	"testing"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/dgraph-io/dgo/y"
	. "github.com/franela/goblin"
)

func TestTestBuilder(t *testing.T) {
	g := Goblin(t)

	g.Describe("testBuilder", func() {
		g.It("should not error when writing", func() {
			b := &testBuilder{}

			n, err := b.Write([]byte{})

			g.Assert(n).
				Equal(0)

			g.Assert(err).
				Equal(error(nil))
		})

		g.It("should error when writing", func() {
			b := &testBuilder{failOn: 1}

			n, err := b.Write([]byte{})

			g.Assert(n).
				Equal(0)

			g.Assert(err).
				Equal(errors.New("WRITE_ERROR"))
		})

		g.It("should not error when stringing", func() {
			b := &testBuilder{}

			n, err := b.Write([]byte{})

			g.Assert(n).
				Equal(0)

			g.Assert(err).
				Equal(error(nil))
		})

		g.It("should error when stringing", func() {
			b := &testBuilder{failOn: 1}

			g.Assert(b.String()).Equal("STRING_ERROR")
		})
	})
}

func TestTestTxn(t *testing.T) {
	g := Goblin(t)

	g.Describe("testTxn", func() {

		g.It("should not error with Query", func() {
			txn := &testTxn{}
			r, err := txn.Query(context.Background(), "")

			g.Assert(r).
				Equal(&api.Response{})

			g.Assert(err).
				Equal(error(nil))
		})

		g.It("should error with Query", func() {
			txn := &testTxn{failOn: 1}
			r, err := txn.Query(context.Background(), "")

			g.Assert(r).
				Equal(&api.Response{})

			g.Assert(err).
				Equal(errors.New("QUERY_ERROR"))
		})

		g.It("should not error with Mutate", func() {
			txn := &testTxn{}
			r, err := txn.Mutate(context.Background(), &api.Mutation{})

			g.Assert(r).
				Equal(&api.Assigned{Uids: map[string]string{"a": "0x1"}})

			g.Assert(err).
				Equal(error(nil))
		})

		g.It("should error with Mutate", func() {
			txn := &testTxn{failOn: 1}
			r, err := txn.Mutate(context.Background(), &api.Mutation{})

			g.Assert(r).
				Equal(&api.Assigned{})

			g.Assert(err).
				Equal(errors.New("MUTATE_ERROR"))
		})

		g.It("should not error with Commit", func() {
			txn := &testTxn{}
			g.Assert(txn.Commit(context.Background())).
				Equal(error(nil))
		})

		g.It("should error with Commit", func() {
			txn := &testTxn{failOn: 1}
			g.Assert(txn.Commit(context.Background())).
				Equal(errors.New("COMMIT_ERROR"))
		})

		g.It("should not error with Discard", func() {
			txn := &testTxn{}
			g.Assert(txn.Discard(context.Background())).
				Equal(error(nil))
		})

		g.It("should error with Discard", func() {
			txn := &testTxn{failOn: 1}
			g.Assert(txn.Discard(context.Background())).
				Equal(errors.New("DISCARD_ERROR"))
		})
	})
}

func TestTestDgraphClient(t *testing.T) {
	g := Goblin(t)

	g.Describe("testDgraphClient", func() {
		dg := &testDgraphClient{}

		g.It("should return response when Login is called", func() {
			r, err := dg.Login(context.Background(), &api.LoginRequest{})
			g.Assert(r).
				Equal(&api.Response{})

			g.Assert(err).
				Equal(error(nil))
		})
		g.It("should return response when Query is called", func() {
			r, err := dg.Query(context.Background(), &api.Request{})
			g.Assert(r).
				Equal(&api.Response{})

			g.Assert(err).
				Equal(error(nil))
		})
		g.It("should return payload when Alter is called", func() {
			r, err := dg.Alter(context.Background(), &api.Operation{})
			g.Assert(r).
				Equal(&api.Payload{})

			g.Assert(err).
				Equal(error(nil))
		})
		g.It("should return context when CommitOrAbort is called", func() {
			r, err := dg.CommitOrAbort(context.Background(), &api.TxnContext{})
			g.Assert(r).
				Equal(&api.TxnContext{})

			g.Assert(err).
				Equal(error(nil))
		})
		g.It("should return version when CheckVersion is called", func() {
			r, err := dg.CheckVersion(context.Background(), &api.Check{})
			g.Assert(r).
				Equal(&api.Version{})

			g.Assert(err).
				Equal(error(nil))
		})
		g.It("should return assigned with a uid when Mutate is called", func() {
			dg.shouldAbort = false
			r, err := dg.Mutate(context.Background(), &api.Mutation{})
			g.Assert(r).
				Equal(&api.Assigned{Uids: map[string]string{"damienstamates": "0x1"}})

			g.Assert(err).
				Equal(error(nil))
		})
		g.It("should return assigned with a uid when Mutate is called", func() {
			dg.shouldAbort = true
			r, err := dg.Mutate(context.Background(), &api.Mutation{})
			g.Assert(r).
				Equal(&api.Assigned{})

			g.Assert(err).
				Equal(y.ErrAborted)
		})
	})
}
