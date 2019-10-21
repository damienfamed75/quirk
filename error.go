package quirk

import (
	"errors"
	"fmt"
)

// Custom errors for Quirk.
var (
	ErrTooManyOperationFields = errors.New("too many fields filled in operation")
	ErrTransactionFailure     = errors.New("transaction failure")
	ErrTooManyResponses       = errors.New("too many responses from query for unique nodes")
	ErrUIDNotFound            = errors.New("couldn't find uid in mutation response")
	ErrNilUID                 = errors.New("*string was nil in response")
)

// QueryError is used for errors revolving around querying.
type QueryError struct {
	ExtErr  error
	Query   string
	Verbose bool
}

// NewQueryError is built for Quirk to fill a custom error type with
// the necessary data to debug an issue quickly, whether or not the
// debug mode is on or off.
func NewQueryError(query string, err error) error {
	return &QueryError{
		Query:  query,
		ExtErr: err,
	}
}

func (e *QueryError) Error() string {
	if e.Verbose {
		return fmt.Sprintf("query error:{%s}:%s", e.Query, e.ExtErr.Error())
	}

	return fmt.Sprintf("query error: %s", e.ExtErr.Error())
}

// Unwrap is used to return any external errors.
// This function was implemented for Go 1.13.
func (e *QueryError) Unwrap() error {
	if e.ExtErr != nil {
		return e.ExtErr
	}

	return nil
}

func (e *QueryError) setVerbose(val bool) *QueryError {
	e.Verbose = val
	return e
}

// MutationError is used when returning a mutation/insertion/update error.
type MutationError struct {
	ExtErr     error
	Identifier string
	RDF        string
	Verbose    bool
	New        bool
}

// NewMutationError is built for Quirk to fill a custom error type with
// the necessary data to debug an issue quickly, whether or not the
// debug mode is on or off.
func NewMutationError(rdf, identifier string, err error) error {
	return &MutationError{
		RDF:        rdf,
		Identifier: identifier,
		ExtErr:     err,
	}
}

func (e *MutationError) Error() string {
	if e.Verbose {
		return fmt.Sprintf("mutation error (%s:%s):{%s}:%s",
			e.Identifier, e.newVal(), e.RDF, e.ExtErr)
	}

	return fmt.Sprintf("mutation error (%s:%s): %s",
		e.Identifier, e.newVal(), e.ExtErr)
}

// Unwrap is used to return any external errors.
// This function was implemented for Go 1.13.
func (e *MutationError) Unwrap() error {
	if e.ExtErr != nil {
		return e.ExtErr
	}

	return nil
}

func (e *MutationError) newVal() string {
	if e.New {
		return "insert"
	}
	return "update"
}

func (e *MutationError) setVerbose(val bool) *MutationError {
	e.Verbose = val
	return e
}

func (e *MutationError) setNew(val bool) *MutationError {
	e.New = val
	return e
}
