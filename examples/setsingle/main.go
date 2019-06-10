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
	SSN    string `quirk:"ssn"`
	Policy string `quirk:"policy"`
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
	c, err := quirk.NewClient(schema)
	if err != nil {
		log.Fatalf("Failed to create Quirk Client [%v]\n", err)
	}

	// Use the quirk client to insert a single node
	// and make sure that if any of the fields are upsert predicates
	// to fail them on transaction and return promptly via the error.
	uidMap, err := c.InsertNode(context.Background(), dg,
		&quirk.Options{SetSingleStruct: &Person{Name: "John", SSN: "126", Policy: "JKL"}},
	)
	if err != nil {
		// If the error is a list of our failed upserts
		// then let's print them out for fun.
		if fUpsert, ok := err.(*quirk.FailedUpsert); ok {
			printFailedUpserts(fUpsert)
		} else {
			log.Fatalf("Error when inserting nodes [%v]\n", err)
		}
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

func printFailedUpserts(fUpsert *quirk.FailedUpsert) {
	log.Printf("FailedUpsertRDF: [%s]\n", fUpsert.GetRDF())
	for _, dat := range fUpsert.GetPredicateValueSlice() {
		log.Printf("FailedPredicateMap: [%s] [%v]\n", dat.Predicate, dat.Value)
	}
}
