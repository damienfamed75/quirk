package quirk

import (
	"errors"
	"testing"

	. "github.com/franela/goblin"
)

func TestError(t *testing.T) {
	g := Goblin(t)

	g.Describe("Error", func() {
		g.It("with external error", func() {
			e := &Error{ExtErr: errors.New("err")}
			g.Assert(e.Error()).
				Equal(":: msg[] external_err[err]")
		})

		g.It("without external error", func() {
			e := &Error{}
			g.Assert(e.Error()).
				Equal(":: msg[]")
		})
	})
}

func TestQueryError(t *testing.T) {
	g := Goblin(t)

	g.Describe("Query error", func() {
		g.It("with external error", func() {
			e := &QueryError{ExtErr: errors.New("err")}
			g.Assert(e.Error()).
				Equal(":query.go: Query[] external_err[err]")
		})

		g.It("with message", func() {
			e := &QueryError{Msg: "msg"}
			g.Assert(e.Error()).
				Equal(":query.go: Msg[msg]")
		})

		g.It("empty", func() {
			e := &QueryError{}
			g.Assert(e.Error()).
				Equal(":query.go: Query[]")
		})
	})
}

func TestTransactionError(t *testing.T) {
	g := Goblin(t)

	g.Describe("Error", func() {
		g.It("with external error", func() {
			e := &TransactionError{ExtErr: errors.New("err")}
			g.Assert(e.Error()).
				Equal(":mutate_single.go: Msg[] RDF[] external_err[err]")
		})

		g.It("without external error", func() {
			e := &TransactionError{}
			g.Assert(e.Error()).
				Equal(":mutate_single.go: Msg[] RDF[]")
		})
	})
}
