package quirk

import (
	"fmt"
)

type Error struct {
	ExtErr   error
	Msg      string
	File     string
	Function string
}

func (e *Error) Error() string {
	if e.ExtErr != nil {
		return fmt.Sprintf("%s:%s: msg[%s] external_err[%v]",
			e.Function, e.File, e.Msg, e.ExtErr,
		)
	}

	return fmt.Sprintf("%s:%s: msg[%s]",
		e.Function, e.File, e.Msg,
	)
}

type FailedUpsert struct {
	PredMap map[string]interface{}
	RDF     string
}

func (e *FailedUpsert) Error() string {
	return "FailedUIDMap: One or more upserts failed. This is not an error."
}

// GetPredicateMap will return the belonging predicate map.
func (e *FailedUpsert) GetPredicateMap() map[string]interface{} {
	return e.PredMap
}

// GetRDF will return the belonging RDF string.
func (e *FailedUpsert) GetRDF() string {
	return e.RDF
}
