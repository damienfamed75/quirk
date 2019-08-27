<p align="center"><img src="https://i.imgur.com/LcLsTX8.png" width="120"></p>
<p align="center"><img src="https://i.imgur.com/L0QU8td.png" width="200"></p>

<p align="center">
<a href="https://godoc.org/github.com/damienfamed75/quirk"><img src="https://godoc.org/github.com/damienfamed75/quirk?status.svg" alt="GoDoc" /></a>
<a href="https://goreportcard.com/report/github.com/damienfamed75/quirk"><img src="https://goreportcard.com/badge/github.com/damienfamed75/quirk" alt="Go Report Card" /></a>
<a href="https://github.com/damienfamed75/quirk/blob/master/LICENSE"><img src="https://img.shields.io/github/license/damienfamed75/quirk.svg" alt="License" /></a>
<a href="https://github.com/damienfamed75/quirk/actions"><img src="https://github.com/damienfamed75/quirk/workflows/Go/badge.svg" /></a>
</p>
<a href="https://codecov.io/gh/damienfamed75/quirk">
  <img src="https://codecov.io/gh/damienfamed75/quirk/branch/master/graph/badge.svg" />
</a>
<p align="center">Quirk is a library used to seemlessly use upsert procedures in Dgraph without going through the hassle yourself.</p>

## Install

Run this command to download the package.

```sh
go get github.com/damienfamed75/quirk
```

## Using quirk

Here is a quick example of using a quirk client to insert a single node.

```go
package main

import (
    "context"
    "fmt"

    "github.com/dgraph-io/dgo/protos/api"
    "github.com/damienfamed75/quirk"
    "github.com/dgraph-io/dgo"
    "google.golang.org/grpc"
)

func main() {
    // Create some data to insert in Dgraph.
    person := struct {
        // These quirk tags are required.
        // It lets the quirk client know that this is the name
        // of the predicate in Dgraph.
        Name   string `quirk:"name"`
        SSN    string `quirk:"ssn,unique"` // Add unique if it should be upserted.
        Policy string `quirk:"policy,unique"`
    }{
        Name:   "Damien",
        SSN:    "123-12-1234",
        Policy: "ABCDAMIEN",
    }

    // Dial with GRPC to Dgraph as usual.
    conn, err := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
    if err != nil {
        fmt.Println(err)
    }
    defer conn.Close()

    // Create the normal dgo client as usual.
    dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

    // Add the schema to Dgraph.
    // Make sure to mark the unique predicates with the "@upsert" directive.
    err = dg.Alter(context.Background(), &api.Operation{
        Schema: `
            name:   string @index(hash) .
            ssn:    string @index(hash) @upsert .
            policy: string @index(hash) @upsert .
        `,
    })
    if err != nil {
        fmt.Println(err)
    }

    // Create a new quirk client.
    q := quirk.NewClient()

    // Insert a single node. If we run this file multiple times you will
    // see that this node is never added twice.
    uids, err := q.InsertNode(context.Background(), dg, &quirk.Operation{
        SetSingleStruct: &person,
    })
    if err != nil {
        fmt.Println(err)
    }

    // Print out the returned node and its uid.
    for n, u := range uids {
        fmt.Printf("UIDMap: name[%s] uid[%s]\n", n, u)
    }
}
```

## Contributing

Go ahead and create a PR or an issue. I'm always open for new contributions.
