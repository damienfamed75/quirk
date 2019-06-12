package quirk

import (
	"reflect"
	"strings"
)

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func reflectMaps(d interface{}) []*predValDat {
	elem := reflect.ValueOf(d).Elem()
	predVal := make([]*predValDat, elem.NumField())

	// loop through elements of struct.
	for i := 0; i < elem.NumField(); i++ {
		var isUpsert bool
		tag, opt := parseTag(reflect.TypeOf(d).Elem().Field(i).Tag.Get("quirk"))

		// store upsert predicates in separate map.
		if opt == tagUnique {
			// If this is an upsert then mark it as such.
			isUpsert = true
		}
		
		// Add the predicate and value to the slice.
		predVal[i] = &predValDat{Predicate: tag, Value: elem.Field(i).Interface(), IsUpsert: isUpsert}
	}

	return predVal
}