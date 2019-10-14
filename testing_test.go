package quirk

import (
	"context"
	"errors"
	"testing"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
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
			_, err := dg.Query(context.Background(), &api.Request{})

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
		// Updated api.Assigned to api.Response as part of dgo major version change.
		// See Issue #17 for more info.
		g.It("should return assigned with a uid when Mutate is called", func() {
			dg.shouldAbort = false
			r, err := dg.Mutate(context.Background(), &api.Mutation{})
			g.Assert(r).
				Equal(&api.Response{Uids: map[string]string{"damienstamates": "0x1"}})

			g.Assert(err).
				Equal(error(nil))
		})
		// Updated api.Assigned to api.Response as part of dgo major version change.
		// See Issue #17 for more info.
		g.It("should return assigned with a uid when Mutate is called", func() {
			dg.shouldAbort = true
			r, err := dg.Mutate(context.Background(), &api.Mutation{})
			g.Assert(r).
				Equal(&api.Response{})

			g.Assert(err).
				Equal(dgo.ErrAborted)
		})
	})
}
