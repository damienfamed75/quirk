package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/damienfamed75/quirk"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

const (
	schema = `
	name: string @index(hash) @upsert .
	age: string .
	`
)

// Person is a way to capture when we need to implement
// a true upsert functionality for when we're dealing with
// social security numbers and policy numbers.
type Person struct {
	Name string `quirk:"name,unique"`
	Age  string `quirk:"age"`
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
	c := quirk.NewClient(quirk.WithLogger(quirk.NewDebugLogger()))

	// In order to insert multiple nodes using the quirk client
	// you must use a slice of interface to as the argument.
	var people []interface{}

	// Opening our testing data.
	f, err := os.Open("data.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		// Create structs from the test data.
		people = append(people, &Person{Name: line[0], Age: line[1]})
	}

	// Use the quirk client to insert multiple nodes at a time
	// all while making sure that any upsert predicates are failed
	// on transaction and returned promptly via the error.

	// We're not storing the uidMap for successful nodes, because
	// the nodes that we have so many nodes that are being inserted
	// and we just want to focus on the failed upserts.
	_, err = c.InsertNode(context.Background(), dg,
		&quirk.Operation{SetMultiStruct: people},
	)
	if err != nil {
		log.Fatalf("Error when inserting nodes [%v]\n", err)
	}
}
