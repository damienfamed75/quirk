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

// Person is a way to capture when we need to implement
// a true upsert functionality for when we're dealing with
// social security numbers and policy numbers.
type Person struct {
	Name   string `quirk:"name"`
	SSN    string `quirk:"ssn,unique"`
	Policy string `quirk:"policy,unique"`
}

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

	// In order to insert multiple nodes using the quirk client
	// you must use a slice of interface to as the argument.
	var people = []interface{}{
		// Damien or George should fail because they share the same Policy.
		&Person{Name: "Damien", SSN: "123", Policy: "ABC"},
		&Person{Name: "George", SSN: "124", Policy: "ABC"},
		// Bahram or Angad should fail because they share the same SSN.
		&Person{Name: "Bahram", SSN: "125", Policy: "DEF"},
		&Person{Name: "Angad", SSN: "125", Policy: "GHI"},
	}

	// Use the quirk client to insert multiple nodes at a time.
	uidMap, err := c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetMultiStruct: people})
	if err != nil {
		log.Fatalf("Error when inserting nodes [%v]\n", err)
	}

	// Finally print out the successful UIDs.
	// The key is typically going to be either your
	// assigned "name" predicate or if you don't have this
	// then it will be an incremented character/s.
	// Note: If you wish to use another predicate beside "name"
	// you may set that when creating the client and using
	// quirk.WithPredicateKey(predicateName string)
	for k, v := range uidMap {
		log.Printf("UIDMap: [%s] [%v:new(%v)]\n", k, v.Value(), v.IsNew())
	}
}
