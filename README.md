<p align="center"><img src="https://i.imgur.com/LcLsTX8.png" width="120"></p>
<p align="center"><img src="https://i.imgur.com/L0QU8td.png" width="200"></p>

<p align="center">
<a href="https://godoc.org/github.com/damienfamed75/quirk"><img src="https://godoc.org/github.com/damienfamed75/quirk?status.svg" alt="GoDoc" /></a>
<a href="https://github.com/damienfamed75/quirk/releases"><img src="https://badgen.net/github/release/damienfamed75/quirk"></a>
<a href="https://goreportcard.com/report/github.com/damienfamed75/quirk"><img src="https://goreportcard.com/badge/github.com/damienfamed75/quirk" alt="Go Report Card" /></a>
<a href="https://github.com/damienfamed75/quirk/blob/master/LICENSE"><img src="https://img.shields.io/github/license/damienfamed75/quirk.svg" alt="License" /></a>
<a href="https://github.com/damienfamed75/quirk/actions"><img src="https://github.com/damienfamed75/quirk/workflows/Build/badge.svg" /></a>
<a href="https://codecov.io/gh/damienfamed75/quirk"><img src="https://codecov.io/gh/damienfamed75/quirk/branch/master/graph/badge.svg"/></a>
</p>

<p align="center">Quirk is a library used to seamlessly use upsert procedures in Dgraph without going through the hassle yourself.</p>

## Some Quick Notes about Quirk v2

With the recent update of Dgraph and dgo, Quirk has been versioned to v2.0.0 and is now using the new module tag `github.com/damienfamed75/quirk/v2`. If you are using any Dgraph version below 1.1.0 and wish to still use quirk then please go ahead and download quirk v1.

Also Quirk will be updating to utilize further functionality added to dgo to optimize mutations as soon as possible. At the moment please enjoy Quirk v2 though in all its glory!

## Install

To download Quirk v2, run this command to download the package.

```sh
go get github.com/damienfamed75/quirk/v2
```

If you wish to download Quirk v1, then run this command instead.

```sh
go get github.com/damienfamed75/quirk
```

### Note if you are using Quirk v1 all imports in Go must not contain /v2

## Using quirk

Here is a quick example of using a quirk client to insert a single node.

```go
package main

import (
    "context"
    "fmt"

    "github.com/damienfamed75/quirk/v2"
    "github.com/dgraph-io/dgo/v2/protos/api"
    "github.com/dgraph-io/dgo/v2"
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
        // As part of Dgraph 1.1.0 you can add types by just using the
        // "dgraph.type" tag using a string and Quirk will handle it correctly.
        Type string `quirk:"dgraph.type"`
    }{
        Name:   "Damien",
        SSN:    "123-12-1234",
        Policy: "ABCDAMIEN",
        Type:   "Person",
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
