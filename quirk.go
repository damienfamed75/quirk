// Package quirk provides the main quirk.Client which is used to insert nodes into
// Dgraph. This client works with multiple types of data and will insert them
// concurrently if using the `multi` Operations found in quirk.Operation.
//
// To get started using the quirk client you must create a Dgraph client using
// Dgraph's official `dgo` package at:
//
// 	https://github.com/dgraph-io/dgo
//
// Once creating a Dgraph client you may begin using the functionality of the quirk
// client. Which currently is just using the `InsertNode` function.
//
// To see examples on how to use the quirk client further then check out the
// 	`examples/`
// directory found in the root of the repository.
// Or you can check out the Godoc examples listed below.
package quirk
