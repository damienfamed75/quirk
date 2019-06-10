package quirk

import (
	"testing"
)

func BenchmarkCreateRDFAuto(b *testing.B) {
	c := &Client{quirkID: aeshash("quirk"), quirkRel: "has"}
	predVal := []*predValDat{
		&predValDat{predicate: "name", value: "damien"},
		&predValDat{predicate: "phone", value: "+1(234)567-8901"},
		&predValDat{predicate: "company", value: "northwestern mutual"},
		&predValDat{predicate: "job", value: "software developer"},
		&predValDat{predicate: "age", value: "19"},
		&predValDat{predicate: "gender", value: "male"},
		&predValDat{predicate: "prodNum", value: "1234-5678-901234"},
		&predValDat{predicate: "favoriteColor", value: "pink"},
		&predValDat{predicate: "favoriteFood", value: "pasta"},
		&predValDat{predicate: "favoriteHobby", value: "drawing"},
	}

	for i := 0; i < b.N; i++ {
		_ = c.createRDF(auto, predVal)
	}
}

func BenchmarkCreateRDFIncrementor(b *testing.B) {
	c := &Client{quirkID: aeshash("quirk"), quirkRel: "has"}
	predVal := []*predValDat{
		&predValDat{predicate: "name", value: "damien"},
		&predValDat{predicate: "phone", value: "+1(234)567-8901"},
		&predValDat{predicate: "company", value: "northwestern mutual"},
		&predValDat{predicate: "job", value: "software developer"},
		&predValDat{predicate: "age", value: "19"},
		&predValDat{predicate: "gender", value: "male"},
		&predValDat{predicate: "prodNum", value: "1234-5678-901234"},
		&predValDat{predicate: "favoriteColor", value: "pink"},
		&predValDat{predicate: "favoriteFood", value: "pasta"},
		&predValDat{predicate: "favoriteHobby", value: "drawing"},
	}

	for i := 0; i < b.N; i++ {
		_ = c.createRDF(incrementor, predVal)
	}
}

func BenchmarkCreateRDFHash(b *testing.B) {
	c := &Client{quirkID: aeshash("quirk"), quirkRel: "has"}
	predVal := []*predValDat{
		&predValDat{predicate: "name", value: "damien"},
		&predValDat{predicate: "phone", value: "+1(234)567-8901"},
		&predValDat{predicate: "company", value: "northwestern mutual"},
		&predValDat{predicate: "job", value: "software developer"},
		&predValDat{predicate: "age", value: "19"},
		&predValDat{predicate: "gender", value: "male"},
		&predValDat{predicate: "prodNum", value: "1234-5678-901234"},
		&predValDat{predicate: "favoriteColor", value: "pink"},
		&predValDat{predicate: "favoriteFood", value: "pasta"},
		&predValDat{predicate: "favoriteHobby", value: "drawing"},
	}

	for i := 0; i < b.N; i++ {
		_ = c.createRDF(hash, predVal)
	}
}

func BenchmarkFmtIDAuto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmtID(auto, 12)
	}
}

func BenchmarkFmtIDOther(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmtID(incrementor, 12)
	}
}
