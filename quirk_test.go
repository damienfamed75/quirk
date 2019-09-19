package quirk_test

import (
	"context"
	"fmt"
	"log"

	"github.com/damienfamed75/quirk"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

func Example_insertSingleNode() {
	// Ignoring error handling for brevity.

	// Dial for Dgraph using grpc.
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	// Create a new Dgraph client for our mutations.
	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	// Alter the schema to be equal to our schema variable.
	dg.Alter(context.Background(), &api.Operation{Schema: `
		name: string @index(hash) .
		ssn: string @index(hash) @upsert .
		policy: string @index(hash) @upsert .
	`})

	// Create the Quirk Client with a debug logger.
	// The debug logger is just for demonstration purposes or for debugging.
	c := quirk.NewClient(quirk.WithLogger(quirk.NewDebugLogger()))

	// Use the quirk client to insert a single node.
	uidMap, _ := c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetSingleStruct: &struct {
			Name   string `quirk:"name"`
			SSN    string `quirk:"ssn,unique"`
			Policy int    `quirk:"policy,unique"`
		}{
			Name:   "Damien",
			SSN:    "126",
			Policy: 61238,
		},
	})

	// Finally print out the UIDs of the nodes that were inserted or updated.
	// The key is going to be either your assigned "name" predicate.
	// Note: If you wish to use another predicate beside "name"
	// you may set that when creating the client and using
	// quirk.WithPredicateKey(predicateName string)
	for k, v := range uidMap {
		log.Printf("UIDMap: [%s] [%s:%v]\n", k, v.Value(), v.IsNew())
	}
}

func ExampleClient_GetPredicateKey() {
	client := quirk.NewClient()

	// Should be "name"
	fmt.Println(client.GetPredicateKey())

	_ = client
}

func ExampleClient_InsertNode_singleStruct() {
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	c := quirk.NewClient()

	// Single struct must be a reference to a struct in order for
	// reflect to not throw a panic.
	data := &struct {
		Name   string `quirk:"name"`
		SSN    string `quirk:"ssn,unique"`
		Policy int    `quirk:"policy,unique"`
	}{
		Name:   "Damien",
		SSN:    "126",
		Policy: 61238,
	}

	c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetSingleStruct: data,
	})
}

func ExampleClient_InsertNode_multiStruct() {
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	c := quirk.NewClient()

	type Person struct {
		Name   string `quirk:"name"`
		SSN    string `quirk:"ssn,unique"`
		Policy int    `quirk:"policy,unique"`
	}

	// Multi structs must be inserted using a slice of interfaces.
	data := []interface{}{
		&Person{
			Name:   "Damien",
			SSN:    "126",
			Policy: 61238,
		}, &Person{
			Name:   "George",
			SSN:    "125",
			Policy: 67234,
		},
	}

	c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetMultiStruct: data,
	})
}

func ExampleClient_InsertNode_singleDupleNode() {
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	c := quirk.NewClient()

	data := &quirk.DupleNode{
		Identifier: "Damien",
		Duples: []quirk.Duple{
			{Predicate: "name", Object: "Damien"},
			{Predicate: "ssn", Object: "126", IsUnique: true},
			{Predicate: "policy", Object: 61238, IsUnique: true},
		},
	}

	c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetSingleDupleNode: data,
	})
}

func ExampleClient_InsertNode_multiDupleNode() {
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	c := quirk.NewClient()

	data := []*quirk.DupleNode{
		&quirk.DupleNode{
			Identifier: "Damien",
			Duples: []quirk.Duple{
				{Predicate: "name", Object: "Damien"},
				{Predicate: "ssn", Object: "126", IsUnique: true},
				{Predicate: "policy", Object: 61238, IsUnique: true},
			},
		},
		&quirk.DupleNode{
			Identifier: "George",
			Duples: []quirk.Duple{
				{Predicate: "name", Object: "George"},
				{Predicate: "ssn", Object: "125", IsUnique: true},
				{Predicate: "policy", Object: 67234, IsUnique: true},
			},
		},
	}

	c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetMultiDupleNode: data,
	})
}

func ExampleClient_InsertNode_stringMap() {
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	c := quirk.NewClient()

	// Maps do not support unique predicates.
	data := make(map[string]string)
	data["name"] = "Damien"
	data["age"] = "19"

	c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetStringMap: data,
	})
}

func ExampleClient_InsertNode_dynamicMap() {
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	defer conn.Close()

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	c := quirk.NewClient()

	// Maps do not support unique predicates.
	// Interface maps (Dynamic maps) support multiple datatypes in Dgraph.
	data := make(map[string]interface{})
	data["name"] = "Damien"
	data["age"] = 19

	c.InsertNode(context.Background(), dg, &quirk.Operation{
		SetDynamicMap: data,
	})
}

func ExampleWithLogger() {
	client := quirk.NewClient(
		quirk.WithLogger(quirk.NewDebugLogger()),
	)

	_ = client
}

func ExampleWithPredicateKey() {
	client := quirk.NewClient(
		quirk.WithPredicateKey("label"),
	)

	// Now that the predicate key is set to "label" the returning UID map
	// will use the predicate value of "label" for the key rather than "name"

	// Should be "label"
	fmt.Println(client.GetPredicateKey())

	_ = client
}

func ExampleWithTemplate() {
	// Using github.com/cheggaaa/pb/v3 for progress bar.
	client := quirk.NewClient(
		quirk.WithTemplate(`{{ "Custom:" }} {{ bar . "[" "-" (cycle . ">" ) " " "]"}} [{{etime . | cyan }}:{{rtime . | cyan }}]`),
	)

	_ = client
}

func ExampleDuple() {
	duple := &quirk.Duple{
		Predicate: "name",
		Object:    "Damien",
		IsUnique:  false,
	}

	_ = duple
}

func ExampleDupleNode() {
	node := &quirk.DupleNode{
		Identifier: "person", // used for the key in the returned UID Map.
		Duples: []quirk.Duple{ // predicate value pairs.
			{Predicate: "name", Object: "Damien"},
			{Predicate: "ssn", Object: "126", IsUnique: true},
			{Predicate: "policy", Object: 61238, IsUnique: true},
		},
	}

	_ = node
}

func ExampleDupleNode_AddDuples() {
	node := &quirk.DupleNode{}

	// Adds new Duples. This doesn't support updating predicate values.
	node.AddDuples(
		quirk.Duple{Predicate: "age", Object: 20},
		quirk.Duple{Predicate: "username", Object: "damienfamed75"},
	)

	_ = node
}

func ExampleDupleNode_Find() {
	node := &quirk.DupleNode{
		Duples: []quirk.Duple{
			{Predicate: "name", Object: "Damien"},
		},
	}

	// Find a Duple stored in the DupleNode.
	name := node.Find("name")

	fmt.Printf("%s: %v unique[%v]\n", name.Predicate, name.Object, name.IsUnique)

	_ = node
}

func ExampleDupleNode_SetOrAdd() {
	node := &quirk.DupleNode{
		Duples: []quirk.Duple{
			{Predicate: "age", Object: 19},
		},
	}

	// Updates the previous valued stored in DupleNode.
	node.SetOrAdd(
		quirk.Duple{Predicate: "age", Object: 20},
	)

	_ = node
}

func ExampleDupleNode_Unique() {
	node := &quirk.DupleNode{
		Identifier: "person", // used for the key in the returned UID Map.
		Duples: []quirk.Duple{ // predicate value pairs.
			{Predicate: "name", Object: "Damien"},
			{Predicate: "ssn", Object: "126", IsUnique: true},
			{Predicate: "policy", Object: 61238, IsUnique: true},
		},
	}

	// returns a slice of all the predicates labeled as unique.
	uniqueDuples := node.Unique()

	// Should be 2.
	fmt.Printf("num of unique nodes [%v]\n", len(uniqueDuples))

	_ = node
}
