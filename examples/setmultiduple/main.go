package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/damienfamed75/quirk"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
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

	// Create some data to insert into Dgraph.
	dupleNodes := []*quirk.DupleNode{
		// Damien or George shuold fail because they share the same Policy.
		{Identifier: "Damien",
			Duples: []quirk.Duple{
				{Predicate: "name", Object: "Damien"},
				{Predicate: "ssn", Object: "123", IsUnique: true},
				{Predicate: "policy", Object: "ABC", IsUnique: true},
			}},
		{Identifier: "George",
			Duples: []quirk.Duple{
				{Predicate: "name", Object: "George"},
				{Predicate: "ssn", Object: "124", IsUnique: true},
				{Predicate: "policy", Object: "ABC", IsUnique: true},
			}},
		// Bahram or Angad should fail because they share the same SSN.
		{Identifier: "Bahram",
			Duples: []quirk.Duple{
				{Predicate: "name", Object: "Bahram"},
				{Predicate: "ssn", Object: "125", IsUnique: true},
				{Predicate: "policy", Object: "DEF", IsUnique: true},
			}},
		{Identifier: "Angad",
			Duples: []quirk.Duple{
				{Predicate: "name", Object: "Angad"},
				{Predicate: "ssn", Object: "125", IsUnique: true},
				{Predicate: "policy", Object: "GHI", IsUnique: true},
			}},
	}

	begin := time.Now()

	// Use the quirk client to insert a single node using a DupleNode struct.
	uidMap, err := c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetMultiDupleNode: dupleNodes})
	if err != nil {
		log.Fatalf("Error when inserting nodes [%v]\n", err)
	}

	log.Printf("Time for insertion [%v]\n", time.Since(begin))

	// Finally print out the successful UIDs.
	// When using DupleNodes the key of the string will be what you put as the
	// node's identifier. If not applicable then the keys will be set to the
	// default, "data"
	for k, v := range uidMap {
		log.Printf("UIDMap: [%s] [%v]\n", k, v)
	}
}
