package quirk

import (
	"testing"
)

func BenchmarkCreateRDFAuto(b *testing.B) {
	c := &Client{quirkID: aeshash("quirk"), quirkRel: "has"}
	predVal := []*PredValDat{
		&PredValDat{Predicate: "name", Value: "damien"},
		&PredValDat{Predicate: "phone", Value: "+1(234)567-8901"},
		&PredValDat{Predicate: "company", Value: "northwestern mutual"},
		&PredValDat{Predicate: "job", Value: "software developer"},
		&PredValDat{Predicate: "age", Value: "19"},
		&PredValDat{Predicate: "gender", Value: "male"},
		&PredValDat{Predicate: "prodNum", Value: "1234-5678-901234"},
		&PredValDat{Predicate: "favoriteColor", Value: "pink"},
		&PredValDat{Predicate: "favoriteFood", Value: "pasta"},
		&PredValDat{Predicate: "favoriteHobby", Value: "drawing"},
	}

	for i := 0; i < b.N; i++ {
		_ = c.createRDF(auto, predVal)
	}
}

func BenchmarkCreateRDFIncrementor(b *testing.B) {
	c := &Client{quirkID: aeshash("quirk"), quirkRel: "has"}
	predVal := []*PredValDat{
		&PredValDat{Predicate: "name", Value: "damien"},
		&PredValDat{Predicate: "phone", Value: "+1(234)567-8901"},
		&PredValDat{Predicate: "company", Value: "northwestern mutual"},
		&PredValDat{Predicate: "job", Value: "software developer"},
		&PredValDat{Predicate: "age", Value: "19"},
		&PredValDat{Predicate: "gender", Value: "male"},
		&PredValDat{Predicate: "prodNum", Value: "1234-5678-901234"},
		&PredValDat{Predicate: "favoriteColor", Value: "pink"},
		&PredValDat{Predicate: "favoriteFood", Value: "pasta"},
		&PredValDat{Predicate: "favoriteHobby", Value: "drawing"},
	}

	for i := 0; i < b.N; i++ {
		_ = c.createRDF(incrementor, predVal)
	}
}

func BenchmarkCreateRDFHash(b *testing.B) {
	c := &Client{quirkID: aeshash("quirk"), quirkRel: "has"}
	predVal := []*PredValDat{
		&PredValDat{Predicate: "name", Value: "damien"},
		&PredValDat{Predicate: "phone", Value: "+1(234)567-8901"},
		&PredValDat{Predicate: "company", Value: "northwestern mutual"},
		&PredValDat{Predicate: "job", Value: "software developer"},
		&PredValDat{Predicate: "age", Value: "19"},
		&PredValDat{Predicate: "gender", Value: "male"},
		&PredValDat{Predicate: "prodNum", Value: "1234-5678-901234"},
		&PredValDat{Predicate: "favoriteColor", Value: "pink"},
		&PredValDat{Predicate: "favoriteFood", Value: "pasta"},
		&PredValDat{Predicate: "favoriteHobby", Value: "drawing"},
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
