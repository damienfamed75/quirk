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
func reflectMaps(d interface{}) predValPairs {
	elem := reflect.ValueOf(d).Elem()
	predVal := make(predValPairs, elem.NumField())

	var (
		tag string     // stores the name of the field in Dgraph.
		opt tagOptions // stores either "" or "unique" if provided in the tags.
	)

	// loop through elements of struct.
	for i := 0; i < elem.NumField(); i++ {
		tag, opt = parseTag(reflect.TypeOf(d).Elem().Field(i).Tag.Get(quirkTag))

		// Add the predicate and value to the slice.
		predVal[i] = &predValDat{
			predicate: tag, // first quirk tag.
			value:     elem.Field(i).Interface(),
			isUnique:  opt == tagUnique, // if the second option is "unique"
		}
	}

	return predVal
}

// checkType will return an XML datatype tag if
// any valid datatypes have are applicable.
func checkType(val interface{}) string {
	switch val.(type) {
	case int64: // int64 gets handled as a normal int.
		return xsInt
	case int32: // int32 gets handled as a normal int.
		return xsInt
	case int16:
		return xsInt
	case int:
		return xsInt
	case bool:
		return xsBool
	case float32: // float32 gets handled as a general float.
		return xsFloat
	case float64: // float64 gets handled as a general float.
		return xsFloat
	}

	return ""
}
