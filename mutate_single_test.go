package quirk

import (
	"testing"
)

func BenchmarkReflectMapsSlice(b *testing.B) {
	var schema = `
	name: string @index(hash) @upsert .
	company: string .
	age: default .
	`
	var damien = struct {
		Name    string `upsert:"name"`
		Company string `upsert:"company"`
		Age     string `upsert:"age"`
	}{
		Name:    "Damien",
		Company: "northwestern mutual",
		Age:     "19",
	}

	c, _ := NewClient(schema)

	for i := 0; i < b.N; i++ {
		_, _ = c.reflectMaps(&damien)
	}
}
