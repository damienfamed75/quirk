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

func (c *Client) createRDF(mode rdfMode, predVal []*PredValDat) string {
	var (
		id     interface{}
		rdf    strings.Builder
		rdfStr = rdfLine
	)

	if mode == hash {
		rdfStr += rdfRel
		quid, qrel := c.quirkID, c.quirkRel

		for _, dat := range predVal {
			id = aeshash(dat.Value.(string))
			fmt.Fprintf(&rdf, rdfStr, id, dat.Predicate, dat.Value, quid, qrel, id)
		}
		return rdf.String()
	}

	if mode == auto {
		rdfStr = rdfAuto
		for _, dat := range predVal {
			if dat.Predicate == c.predicateKey {
				id = dat.Value
			}
		}
		if id == nil {
			// For every new RDF this incrementor should be reset.
			resetIncrementor(&charIncrementor)
			id = increment(mode)
		}
	} else {
		id = increment(mode)
	}

	// create RDF lines in string.
	for _, dat := range predVal {
		fmt.Fprintf(&rdf, rdfStr, id, dat.Predicate, dat.Value)
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
