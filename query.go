package quirk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// queryUID builds a query, executes it to find any unique marked predicates.
// Then sifts through the response to see if any UIDs were returned.
// Finally returns either a uid or an empty string if nothing was found.
func queryUID(ctx context.Context, txn dgraphTxn, b builder, dat predValPairs) (string, error) {
	defer b.Reset() // reset the strings builder.

	var decode queryDecode // where the decoded query response will be stored.

	// execute a query to find any UIDs that are existing for unique fields.
	if err := executeQuery(ctx, txn, b, dat, &decode); err != nil {
		return "", err
	}

	// find if the decoded query contains any UIDs.
	return findDecodedUID(decode)
}

// findDecodeUID loops through the decoded response and tries to find
// any UIDs. If any are found then they are returned. Otherwise
// this function will just return an empty string.
func findDecodedUID(decode queryDecode) (string, error) {
	for _, v := range decode {
		if len(v) > 1 { // if there are too many responses.
			return "", &QueryError{
				Msg: msgTooManyResponses, Function: "findDecodedUID"}
		}
		if len(v) == 1 {
			if v[0].UID == nil { // if the *string is nil.
				return "", &QueryError{
					Msg: msgNilUID, Function: "findDecodedUID"}
			}
			return *v[0].UID, nil
		}
	}

	return "", nil
}

// executeQuery calls to create a Query based on the unique predicates sent in
// and then executes it using the given transaction.
func executeQuery(ctx context.Context, txn dgraphTxn, b builder,
	dat predValPairs, decode *queryDecode) error {
	// Write a query that finds the marked unique predicates to the builder.
	if err := createQuery(b, dat); err != nil {
		return err
	}

	// If the query is empty then return nil right away.
	if b.String() == emptyQuery {
		return nil
	}

	// Execute the query with the given transaction.
	resp, err := txn.Query(ctx, b.String())
	if err != nil {
		return &QueryError{
			Function: "executeQuery",
			Query:    b.String(),
			ExtErr:   err,
		}
	}

	// Unmarshal the response into the given decode object.
	if err = json.Unmarshal(resp.GetJson(), decode); err != nil {
		return &QueryError{
			Function: "executeQuery",
			ExtErr:   err,
		}
	}

	return nil
}

// createQuery will loop through the unique predicates and write the query
// to the given io.Writer.
func createQuery(b io.Writer, dat predValPairs) error {
	if _, err := b.Write([]byte{'{'}); err != nil {
		return err
	}

	// Loop through and add a new function per unique predicate.
	for _, d := range dat.unique() {
		_, err := fmt.Fprintf(b, queryfunc, "find"+d.predicate, d.predicate, d.value)
		if err != nil { // returns quirk.Error for predicate and value context.
			return &QueryError{
				ExtErr:   err,
				Msg:      fmt.Sprintf(msgBuilderWriting, d.predicate, d.value),
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
