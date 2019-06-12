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
	username: string @index(hash) @upsert .
	company: string @index(hash) .
	website: string @index(hash) @upsert .
	`
)

// Profile is our way of showing off that we don't
// need a "name" predicate for custom keys in the returned
// uid map. You may set the key to be any predicate you wish.
type Profile struct {
	Username string `quirk:"username,unique"`
	Company  string `quirk:"company"`
	Website  string `quirk:"website,unique"`
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

	// Create the Quirk Client with our schema and WithPredicateKey
	// so we can set a custom key to use when creating the UID map result.
	// Note: If this function is not added when creating the Quirk client
	// the default predicate key will be set to "name"
	c := quirk.NewClient(quirk.WithPredicateKey("username"))

	// In order to insert multiple nodes using the quirk client
	// you must use a slice of interface to as the argument.
	var profiles []interface{}

	// One of the two profiles will fail because they both share the same username.
	profiles = append(profiles, &Profile{Username: "damienstamates", Company: "NM", Website: "northwesternmutual.com"})
	profiles = append(profiles, &Profile{Username: "barum", Company: "NM", Website: "northwesternmutual.com"})
	profiles = append(profiles, &Profile{Username: "gevuong", Company: "NM", Website: "northwesternmutual.com"})
	profiles = append(profiles, &Profile{Username: "damienstamates", Company: "FOXCONN", Website: "foxconn.com"})
	profiles = append(profiles, &Profile{Username: "angad", Company: "NM", Website: "northwesternmutual.com"})
	profiles = append(profiles, &Profile{Username: "cyberninja89", Company: "NTT", Website: "nttdata.com"})
	profiles = append(profiles, &Profile{Username: "solarlune", Company: "N/A", Website: "solarlune.com"})
	profiles = append(profiles, &Profile{Username: "happycow77", Company: "SCHUBERT", Website: "shuberthartfordtheater.com"})
	profiles = append(profiles, &Profile{Username: "cyberninja89", Company: "FOXCONN", Website: "foxconn.com"})
	profiles = append(profiles, &Profile{Username: "barum", Company: "FOXCONN", Website: "foxconn.com"})

	// Use the quirk client to insert multiple nodes at a time
	// all while making sure that any upsert predicates are failed
	// on transaction and returned promptly via the error.
	uidMap, err := c.InsertNode(context.Background(), dg,
		&quirk.Operation{SetMultiStruct: profiles},
	)
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
		log.Printf("UIDMap: [%s] [%s]\n", k, v)
	}
}