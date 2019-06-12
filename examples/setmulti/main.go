package main

import (
	"context"
	"log"

	"github.com/damienfamed75/quirk"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

const (
	schema = `
	name: string @index(hash) .
	ssn: string @index(hash) @upsert .
	policy: string @index(hash) @upsert .
	`
)

// Person is a way to capture when we need to implement
// a true upsert functionality for when we're dealing with
// social security numbers and policy numbers.
type Person struct {
	Name   string `quirk:"name"`
	SSN    string `quirk:"ssn,unique"`
	Policy string `quirk:"policy,unique"`
}

func main() {
	// Dial for Dgraph using grpc.
	conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed when dialing grpc [%v]", err)
	}
	defer conn.Close()

	// Create a new Dgraph client for our mutations.
	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	// Drop all pre-existing data in the graph.
	err = dg.Alter(context.Background(), &api.Operation{DropAll: true})
	if err != nil {
		log.Fatalf("Alteration error with DropAll [%v]\n", err)
	}

	// Alter the schema to be equal to our schema variable.
	err = dg.Alter(context.Background(), &api.Operation{Schema: schema})
	if err != nil {
		log.Fatalf("Alteration error with setting schema [%v]\n", err)
	}

	// Create the Quirk Client with our schema.
	// The schema is read and processed so the client knows
	// which predicates use the @upsert directive.
	c := quirk.NewClient()

	// In order to insert multiple nodes using the quirk client
	// you must use a slice of interface to as the argument.
	var people []interface{}

	// Damien or George should fail because they share the same Policy.
	people = append(people, &Person{Name: "Damien", SSN: "123", Policy: "ABC"})
	people = append(people, &Person{Name: "George", SSN: "124", Policy: "ABC"})

	// Bahram or Angad should fail because they share the same SSN.
	people = append(people, &Person{Name: "Bahram", SSN: "125", Policy: "DEF"})
	people = append(people, &Person{Name: "Angad", SSN: "125", Policy: "GHI"})

	// Use the quirk client to insert multiple nodes at a time
	// all while making sure that any upsert predicates are failed
	// on transaction and returned promptly via the error.
	uidMap, err := c.InsertNode(context.Background(), dg,
		&quirk.Operation{SetMultiStruct: people},
	)
	if err != nil {
		// If the error is a list of our failed upserts
		// then let's print them out for fun.
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
		log.Printf("UIDMap: [%s] [%s]\n", k, v)
	}
}