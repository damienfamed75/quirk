// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quirk

import (
	"reflect"
	"strings"
)

// Credit: The Go Authors @ "encoding/json"
// parseTag splits a struct field's quirk tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// reflectMaps takes an interface which is suspected to be given as a struct.
// Takes the fields and parses through the tags and adds the fieldname, and
// whether it is a unique tag to the returning predicate:value data slice.
func (c *Client) reflectMaps(d interface{}) *DupleNode {
	elem := reflect.ValueOf(d).Elem()
	numFields := elem.NumField()
	duple := &DupleNode{
		// Duples: make([]Duple, numFields),
	}

	var (
		tag string     // stores the name of the field in Dgraph.
		opt tagOptions // stores either "" or "unique" if provided in the tags.
	)

	// loop through elements of struct.
	for i := 0; i < numFields; i++ {
		tag, opt = parseTag(reflect.TypeOf(d).Elem().Field(i).Tag.Get(quirkTag))
		if tag == "" {
			continue
		}

		// Add the predicate and value to the slice.
		// duple.Duples[i] = Duple{
		duple.Duples = append(duple.Duples, Duple{
			Predicate: tag, // first quirk tag.
			Object:    elem.Field(i).Interface(),
			IsUnique:  opt == tagUnique, // if the second option is "unique"
			// dataType:  checkType(elem.Field(i).Interface()),
		})

		if tag == c.predicateKey {
			duple.Identifier = duple.Duples[i].Object.(string)
		}
	}

	return duple
}

// dynamicMapToPredValPairs may seem like a duplicate of mapToPredValPairs
// but this one uses a map[string]interface{} instead of a map[string]string
// which in Go there is currently no way to have them be interchangeable.
func (c *Client) dynamicMapToPredValPairs(d map[string]interface{}) *DupleNode {
	duple := &DupleNode{
		Duples: make([]Duple, len(d)),
	}

	var i int // counter for the predVal slice.
	var val interface{}

	// loop through elements of map.
	for k, v := range d {

		dType := checkType(v)
		if dType == xsByte {
			val = string(v.([]byte))
			dType = ""
		} else {
			val = v
		}

		duple.Duples[i] = Duple{
			Predicate: k,
			Object:    val,
			dataType:  dType,
		}

		if k == c.predicateKey {
			duple.Identifier = v.(string)
		}

		i++
	}

	return duple
}

// mapToPredValPairs may seem like a duplicate of dynamicMapToPredValPairs
// but this one uses a map[string]string instead of a map[string]interface{}
// which in Go there is currently no way to have them be interchangeable.
func (c *Client) mapToPredValPairs(d map[string]string) *DupleNode {
	duple := &DupleNode{
		Duples: make([]Duple, len(d)),
	}

	var i int // counter for the predVal slice.

	// loop through elements of map.
	for k, v := range d {
		duple.Duples[i] = Duple{
			Predicate: k,
			Object:    v,
		}

		if k == c.predicateKey {
			duple.Identifier = v
		}

		i++
	}

	return duple
}

// checkType will return an XML datatype tag if
// any valid datatypes have are applicable.
func checkType(val interface{}) string {
	switch val.(type) {
	case int64: // int64 gets handled as a normal int.
		return xsInt
	case int32: // int32 gets handled as a normal int.
		return xsInt
	case int16: // int16 gets handled as normal int.
		return xsInt
	case int8: // int8 gets handled as normal int.
		return xsInt
	case int:
		return xsInt
	case bool:
		return xsBool
	case float32: // float32 gets handled as a general float.
		return xsFloat
	case float64: // float64 gets handled as a general float.
		return xsFloat
	case []byte:
		return xsByte
	}

	return ""
}
