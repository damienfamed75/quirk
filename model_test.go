package quirk

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestFind(t *testing.T) {
	g := Goblin(t)
	d := testPredValCorrect.Duples[5]

	g.Describe("Find reference", func() {
		g.It("Found predicate", func() {
			g.Assert(testPredValCorrect.Find("testByteSlice")).
				Equal(&d)
		})

		g.It("Predicate not found", func() {
			g.Assert(testPredValCorrect.Find("Unknown")).
				Equal((*Duple)(nil))
		})
	})
}

func TestSetOrAdd(t *testing.T) {
	g := Goblin(t)
	d := testPredValCorrect

	g.Describe("", func() {
		g.It("Found pre existing duple", func() {
			foundDuple := Duple{Predicate: "username",
				Object:   "eviiviviana",
				IsUnique: true}
			g.Assert(d.SetOrAdd(foundDuple)).
				Equal(testPredValCorrect)

			if testPredValCorrect.Duples[0].Object != "eviiviviana" {
				g.Fail(error(nil))
			}
			//revert to regular testing data
			testPredValCorrect.Duples[0].Object = testPersonCorrect.Username
		})

		g.It("Duple doesn't exist", func() {
			unfoundDuple := Duple{Predicate: "ssn",
				Object:   "1111",
				IsUnique: true}
			g.Assert(d.SetOrAdd(unfoundDuple)).
				Equal(testPredValCorrect)

			if d.Find("ssn") == (*Duple)(nil) {
				g.Fail(error(nil))
			}
		})
	})
}
func TestAddDuples(t *testing.T) {

}
func TestIsNew(t *testing.T) {
	g := Goblin(t)

	g.Describe("New UID", func() {
		g.It("return bool", func() {
			u := UID{uid: "test", isNew: true}

			g.Assert(u.IsNew()).
				Equal(true)
		})
	})
}
