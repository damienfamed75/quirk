package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

//go:noinline
func BenchmarkMain(b *testing.B) {
	ctx := context.Background()
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())

	dgraph := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	var txn *dgo.Txn
	for i := 0; i <= b.N; i++ {
		txn = dgraph.NewReadOnlyTxn()

		txn.BestEffort()
		_, _ = txn.Query(ctx, `
		{
			classy(func: eq(name, "testNode")) {
				first
				last
			}
		}
		`)
	}
}

//go:noinline
func BenchmarkMainCache(b *testing.B) {
	ctx := context.Background()
	conn, _ := grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())

	dgraph := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	txn := dgraph.NewTxn()

	res, _ := txn.Query(ctx, `
	{
		classy(func: eq(name, "testNode")) {
			first
			last
		}
	}
	`)

	var nodes struct {
		Res []map[string]interface{} `json:"classy"`
	}

	json.Unmarshal(res.GetJson(), &nodes)

	for i := 0; i <= b.N; i++ {
		_ = fmt.Sprintf("%s", nodes.Res[0]["first"])
		_ = fmt.Sprintf("%s", nodes.Res[0]["last"])
	}
}
