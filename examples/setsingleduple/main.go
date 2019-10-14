package main

import (
	"context"
	"flag"
	"log"

	"github.com/damienfamed75/quirk/v2"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

var drop bool

// Although a little more complex, using DupleNode will save time
// when inserting nodes using the Quirk client, because then the quirk client
// doesn't have to loop through your data as much.

func main() {
	flag.BoolVar(&drop, "d", false, "Drop-All before running example.")
	flag.Parse()

	// Dial for Dgraph using grpc.
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed when dialing grpc [%v]", err)
	}
	defer conn.Close()

	// Create a new Dgraph client for our mutations.
	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	// Drop all pre-existing data in the graph.
	if drop {
		err = dg.Alter(context.Background(), &api.Operation{DropAll: true})
		if err != nil {
			log.Fatalf("Alteration error with DropAll [%v]\n", err)
		}
	}

	// Alter the schema to be equal to our schema variable.
	err = dg.Alter(context.Background(), &api.Operation{Schema: `
		name: string @index(hash) .
		ssn: string @index(hash) @upsert .
		policy: string @index(hash) @upsert .
	`})
	if err != nil {
		log.Fatalf("Alteration error with setting schema [%v]\n", err)
	}

	// Create the Quirk Client with a debug logger.
	// The debug logger is just for demonstration purposes or for debugging.
	c := quirk.NewClient(quirk.WithLogger(quirk.NewDebugLogger()))
	if err != nil {
		log.Fatalf("Failed to create Quirk Client [%v]\n", err)
	}

	dupleNode := &quirk.DupleNode{
		Identifier: "John",
		Duples: []quirk.Duple{
			{Predicate: "name", Object: "John"},
			{Predicate: "ssn", Object: "126", IsUnique: true},
			{Predicate: "policy", Object: "JKL", IsUnique: true},
		},
	}

	// Use the quirk client to insert a single node using a DupleNode struct.
	uidMap, err := c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetSingleDupleNode: dupleNode})
	if err != nil {
		log.Fatalf("Error when inserting nodes [%v]\n", err)
	}

	// Finally print out the successful UIDs.
	// When using DupleNodes the key of the string will be what you put as the
	// node's identifier. If not applicable then the keys will be set to the
	// default, "data"
	for k, v := range uidMap {
		log.Printf("UIDMap: [%s] [%v]\n", k, v)
	}
}
