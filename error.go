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

// QueryError is used for functions in the query.go file.
type QueryError struct {
	ExtErr   error
	Msg      string
	Function string
	Query    string
}

func (e *QueryError) Error() (res string) {
	switch {
	case e.ExtErr != nil:
		res = fmt.Sprintf("%s:query.go: Query[%s] external_err[%v]",
			e.Function, e.Query, e.ExtErr,
		)
	case e.Msg != "":
		res = fmt.Sprintf("%s:query.go: Msg[%s]",
			e.Function, e.Msg,
		)
	default:
		res = fmt.Sprintf("%s:query.go: Query[%s]",
			e.Function, e.Query,
		)
	}

	return
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
