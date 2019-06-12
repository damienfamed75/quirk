package quirk

import (
	"io"
	"fmt"
	"context"
	"encoding/json"
)

func queryUID(ctx context.Context, txn dgraphTxn, b builder, dat []*predValDat) (string, error) {
	defer b.Reset()

	var decode queryDecode

	if err := executeQueryForUnique(ctx, txn, b, dat, &decode); err != nil {
		return "", err
	}
	
	return assignUIDFromDecode(decode)
}

func assignUIDFromDecode(decode queryDecode) (string, error) {
	for _, v := range decode {
		if len(v) <= 1 != true {
			return "", &Error{Msg: "INVALID LEN"}
		}
		if len(v) == 1 {
			if v[0].UID == nil {
				return "", &Error{Msg: "UID NIL"}
			}
			return *v[0].UID, nil
		}
	}

	return "", nil
}

func createQuery(b io.Writer, dat []*predValDat) error {
	if _, err := b.Write([]byte{'{'}); err != nil {
		return err
	}

	// Loop through and add a new function per unique predicate.
	for _, d := range dat {
		if d.IsUpsert {
			_, err := fmt.Fprintf(b, queryfunc, "f"+d.Predicate, d.Predicate, d.Value)
			if err != nil {
				return &Error{
					ExtErr:   err,
					Msg:      fmt.Sprintf("predicate[%#v] value[%#v]", d.Predicate, d.Value),
					File:     "query.go",
					Function: "createQuery",
				}
			}
		}
	}

	// End the query.
	if _, err := b.Write([]byte{'}'}); err != nil {
		return err
	}

	return nil
}

func executeQueryForUnique(ctx context.Context, txn dgraphTxn, b builder,
	dat []*predValDat, decode *queryDecode) error {
	if err := createQuery(b, dat); err != nil {
		return err
	}

	if b.String() == emptyQuery {
		return nil
	}

	resp, err := txn.Query(ctx, b.String())
	if err != nil {
		return &QueryError{
			ExtErr: err, Query: b.String(),
			File:     "query.go",
			Function: "executeQueryForUnique",
		}
	}

	if err = json.Unmarshal(resp.GetJson(), decode); err != nil {
		return err
	}

	return nil
}
