package quirk

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func TestNewQueryError(t *testing.T) {
	g := Goblin(t)

	g.Describe("Query error", func() {
		g.It("with non-verbose logging", func() {
			e := NewQueryError(_emptyQuery, ErrNilUID)
			g.Assert(e.Error()).
				Equal(fmt.Sprintf("query error:%s",
					ErrNilUID.Error()))
		})
		g.It("with verbose logging", func() {
			e := NewQueryError(_emptyQuery, ErrNilUID)
			g.Assert(e.Error()).
				Equal(fmt.Sprintf("query error:{%s}:%s",
					_emptyQuery, ErrNilUID.Error()))
		})
	})
}
