package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"
	"time"

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
	Name string `quirk:"name,unique"`
	Age  string `quirk:"age"`
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
		name: string @index(hash) @upsert .
		age: string .
	`})
	if err != nil {
		log.Fatalf("Alteration error with setting schema [%v]\n", err)
	}

	// Create the Quirk Client with a debug logger.
	// The debug logger is just for demonstration purposes or for debugging.
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

	begin := time.Now()

	// Use the quirk client to insert multiple nodes at a time
	// all while making sure that any upsert predicates are failed
	// on transaction and returned promptly via the error.
	uidMap, err := c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetMultiStruct: people})
	if err != nil {
		log.Fatalf("Error when inserting nodes [%v]\n", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ',
		tabwriter.AlignRight|tabwriter.Debug)
	var count int

	for k, v := range uidMap {
		count++
		if count%3 == 0 {
			fmt.Fprintf(w, "%s\t%v\t\n", k, v.Value())
		} else {
			fmt.Fprintf(w, "%s\t%v\t", k, v.Value())
		}
	}
	w.Flush()

	fmt.Printf("Inserted nodes in [%v]\n", time.Since(begin))
}
