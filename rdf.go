package quirk

import (
	"fmt"
	"hash/fnv"
)

const (
	rdfRel  = "<%v> <%s> <%v> .\n"
	rdfLine = "<%v> <%s> %#v .\n"
	rdfAuto = "_:%v <%s> %#v .\n"
)

func (c *Client) createRDF(mode rdfMode, predValMap map[string]interface{}) []string {
	var (
		rdf    []string
		rdfStr = rdfLine
	)

	if mode == auto {
		rdfStr = rdfAuto
		// For every new RDF this incrementor should be reset.
		resetIncrementor(&charIncrementor)
	}

	var id interface{}
	if name, ok := predValMap[c.predicateKey]; ok && mode == auto {
		id = name
	} else {
		id = fmtID(mode, increment(mode))
	}

	// create RDF lines in string.
	for k, v := range predValMap {
		if mode == hash {
			id = fingerprint(v.(string))
			rdf = append(rdf, fmt.Sprintf(rdfStr+rdfRel, id, k, v,
				fingerprint(c.quirkName), c.quirkRel, id))
		} else {
			rdf = append(rdf, fmt.Sprintf(rdfStr, id, k, v))
		}
	}

	// use incrementor, hash, or auto as UID.
	return rdf
}

func fingerprint(val string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(val))
	return h.Sum64()
}

func fmtID(mode rdfMode, id uint64) interface{} {
	if mode == auto {
		return string(id)
	}
	return id
}
