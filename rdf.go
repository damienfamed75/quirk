package quirk

import (
	"fmt"
	"strings"
)

const (
	rdfRel  = "<%d> <%s> <%d> .\n"
	rdfLine = "<%v> <%s> %#v .\n"
	rdfAuto = "_:%v <%s> %#v .\n"
)

func (c *Client) createRDF(mode rdfMode, predVal []*predValDat) string {
	var (
		id     interface{}
		rdf    strings.Builder
		rdfStr = rdfLine
	)

	if mode == hash {
		rdfStr += rdfRel
		quid, qrel := c.quirkID, c.quirkRel

		for _, dat := range predVal {
			id = aeshash(dat.value.(string))
			fmt.Fprintf(&rdf, rdfStr, id, dat.predicate, dat.predicate, quid, qrel, id)
		}
		return rdf.String()
	}

	if mode == auto {
		rdfStr = rdfAuto
		// For every new RDF this incrementor should be reset.
		resetIncrementor(&charIncrementor)
		for _, dat := range predVal {
			if dat.predicate == c.predicateKey {
				id = dat.value
			}
		}
	} else {
		id = increment(mode)
	}

	// create RDF lines in string.
	for _, dat := range predVal {
		fmt.Fprintf(&rdf, rdfStr, id, dat.predicate, dat.value)
	}

	// use incrementor, hash, or auto as UID.
	return rdf.String()
}

func fmtID(mode rdfMode, id uint64) interface{} {
	if mode == auto {
		return string(id)
	}
	return id
}
