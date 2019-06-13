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

type QueryError struct {
	ExtErr   error
	Msg      string
	File     string
	Function string
	Query    string
}

func (e *QueryError) Error() string {
	if e.ExtErr != nil {
		return fmt.Sprintf("%s:%s: Query[%s] external_err[%v]",
			e.Function, e.File, e.Query, e.ExtErr,
		)
	}

	return fmt.Sprintf("%s:%s: Query[%s]",
		e.Function, e.File, e.Query,
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
