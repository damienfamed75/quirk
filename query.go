package quirk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func queryUID(ctx context.Context, txn dgraphTxn, b builder, dat predValPairs) (string, error) {
	defer b.Reset()

	var decode queryDecode

	if err := executeQuery(ctx, txn, b, dat, &decode); err != nil {
		return "", err
	}

	return findDecodedUID(decode)
}

func findDecodedUID(decode queryDecode) (string, error) {
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

func executeQuery(ctx context.Context, txn dgraphTxn, b builder,
	dat predValPairs, decode *queryDecode) error {
	if err := createQuery(b, dat); err != nil {
		return err
	}

	if b.String() == emptyQuery {
		return nil
	}

	resp, err := txn.Query(ctx, b.String())
	if err != nil {
		return &QueryError{
			File:     "query.go",
			Function: "executeQuery",
			Msg:      msgQueryingUnique,
			Query:    b.String(),
			ExtErr:   err,
		}
	}

	if err = json.Unmarshal(resp.GetJson(), decode); err != nil {
		return err
	}

	return nil
}

func createQuery(b io.Writer, dat predValPairs) error {
	if _, err := b.Write([]byte{'{'}); err != nil {
		return err
	}

	// Loop through and add a new function per unique predicate.
	for _, d := range dat.unique() {
		_, err := fmt.Fprintf(b, queryfunc, "find"+d.predicate, d.predicate, d.value)
		if err != nil { // returns quirk.Error for predicate and value context.
			return &Error{
				ExtErr:   err,
				Msg:      fmt.Sprintf(msgBuilderWriting, d.predicate, d.value),
				File:     "query.go",
				Function: "createQuery",
			}
		}
	}

	// End the query.
	if _, err := b.Write([]byte{'}'}); err != nil {
		return err
	}

	return nil
}
