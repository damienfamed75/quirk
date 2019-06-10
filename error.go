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

type TransactionError struct {
	ExtErr   error
	Msg      string
	File     string
	Function string
	RDF      string
}

func (e *TransactionError) Error() string {
	if e.ExtErr != nil {
		return fmt.Sprintf("%s:%s: RDF[%s] external_err[%v]",
			e.Function, e.File, e.RDF, e.ExtErr,
		)
	}

	return fmt.Sprintf("%s:%s: RDF[%s]",
		e.Function, e.File, e.RDF,
	)
}

type FailedUpserts struct {
	Upserts []*LoneUpsertError
}

func (e *FailedUpserts) Error() string {
	return "FailedUpserts: One or more upserts failed. This is not an error."
}

func (e *FailedUpserts) Len() int {
	return len(e.Upserts)
}

func (e *FailedUpserts) append(err ...*LoneUpsertError) {
	e.Upserts = append(e.Upserts, err...)
}

type LoneUpsertError struct {
	PredVals []*PredValDat
	RDF      string
}

func (e *LoneUpsertError) Error() string {
	return "LoneUpsertError: This upsert failed. This is not an error."
}

// GetPredicateValueSlice will return the belonging predicate slice.
func (e *LoneUpsertError) GetPredicateValueSlice() []*PredValDat {
	return e.PredVals
}

// GetRDF will return the belonging RDF string.
func (e *LoneUpsertError) GetRDF() string {
	return e.RDF
}
