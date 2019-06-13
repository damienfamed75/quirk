package quirk

import (
	"testing"

	. "github.com/franela/goblin"
)

func BenchmarkReflectMaps(b *testing.B) {
	data := testPersonCorrect // from testing.go

	for i := 0; i < b.N; i++ {
		_ = reflectMaps(&data)
	}
}

func BenchmarkParseTag(b *testing.B) {
	data := "1,2"

	for i := 0; i < b.N; i++ {
		_, _ = parseTag(data)
	}
}

func TestReflectMaps(t *testing.T) {
	g := Goblin(t)

	g.Describe("Reflect Maps", func() {
		g.It("Should equal testPredValCorrect", func(done Done) {
			go func() {
				g.Assert(reflectMaps(&testPersonCorrect)).
					Equal(testPredValCorrect)
				done()
			}()
		})

		g.It("Should equal testPredValInvalid", func(done Done) {
			go func() {
				g.Assert(reflectMaps(&testPersonInvalid)).
					Equal(testPredValInvalid)
				done()
			}()
		})
	})
}

func TestParseTag(t *testing.T) {
	g := Goblin(t)

	g.Describe("Parsed Tags", func() {
		g.It("Should equal \"1\", \"2\"", func(done Done) {
			go func() {
				a, b := parseTag("1,2")
				g.Assert(a).
					Equal("1")
				g.Assert(b).
					Equal(tagOptions("2"))

				done()
			}()
		})

		g.It("Should equal \"1\", \"\"", func(done Done) {
			go func() {
				a, b := parseTag("1,")
				g.Assert(a).
					Equal("1")
				g.Assert(b).
					Equal(tagOptions(""))

				done()
			}()
		})
	})
}
