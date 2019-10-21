package quirk

import (
	"testing"

	. "github.com/franela/goblin"
)

func BenchmarkReflectMaps(b *testing.B) {
	c := NewClient()
	data := testPersonCorrect // from testing.go

	for i := 0; i < b.N; i++ {
		_ = c.reflectMaps(&data)
	}
}

func BenchmarkParseTag(b *testing.B) {
	data := "1,2"

	for i := 0; i < b.N; i++ {
		_, _ = parseTag(data)
	}
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

func TestCheckType(t *testing.T) {
	g := Goblin(t)

	g.Describe("XML Datatype", func() {
		g.It("Should be an int with int", func() {
			var a int
			b := checkType(a)

			g.Assert(b).
				Equal(_xsInt)
		})
		g.It("Should be an int with int64", func() {
			var a int64
			b := checkType(a)

			g.Assert(b).
				Equal(_xsInt)
		})
		g.It("Should be an int with int32", func() {
			var a int32
			b := checkType(a)

			g.Assert(b).
				Equal(_xsInt)
		})
		g.It("Should be an int with int16", func() {
			var a int16
			b := checkType(a)

			g.Assert(b).
				Equal(_xsInt)
		})
		g.It("Should be an int with int8", func() {
			var a int8
			b := checkType(a)

			g.Assert(b).
				Equal(_xsInt)
		})
		g.It("Should be bool with bool", func() {
			var a bool
			b := checkType(a)

			g.Assert(b).
				Equal(_xsBool)
		})
		g.It("Should be float with float32", func() {
			var a float32
			b := checkType(a)

			g.Assert(b).
				Equal(_xsFloat)
		})
		g.It("Should be float with float64", func() {
			var a float64
			b := checkType(a)

			g.Assert(b).
				Equal(_xsFloat)
		})
		g.It("Should be empty with string", func() {
			var a string
			b := checkType(a)

			g.Assert(b).
				Equal("")
		})
		g.It("Should be byte with byte slice", func() {
			var a []byte
			b := checkType(a)

			g.Assert(b).
				Equal(_xsString)
		})
	})
}
