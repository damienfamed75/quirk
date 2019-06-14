package quirk

import (
	"context"
	"errors"

	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// testBuilder is used to mock out a strings.Builder
// and includes a failOn int to specify when the builder
// should return an error back when being used.
type testBuilder struct {
	useCount     int
	failOn       int
	stringOutput string
}

func (t *testBuilder) Write([]byte) (int, error) {
	t.useCount++
	if t.useCount == t.failOn {
		return 0, errors.New("WRITE_ERROR")
	}
	return 0, nil
}

func (t *testBuilder) String() string {
	t.useCount++
	if t.useCount == t.failOn {
		return "STRING_ERROR"
	}
	return t.stringOutput
}

func (*testBuilder) Reset() {}

type testTxn struct {
	useCount   int
	failOn     int
	jsonOutput []byte
}

func (t *testTxn) Query(context.Context, string) (*api.Response, error) {
	t.useCount++
	if t.useCount == t.failOn {
		return &api.Response{}, errors.New("QUERY_ERROR")
	}
	return &api.Response{
		Json: []byte(t.jsonOutput)}, nil
}

func (t *testTxn) Mutate(context.Context, *api.Mutation) (*api.Assigned, error) {
	t.useCount++
	if t.useCount == t.failOn {
		return &api.Assigned{}, errors.New("MUTATE_ERROR")
	}
	return &api.Assigned{Uids: map[string]string{"a": "0x1"}}, nil
}

func (t *testTxn) Commit(context.Context) error {
	t.useCount++
	if t.useCount == t.failOn {
		return errors.New("COMMIT_ERROR")
	}
	return nil
}

func (t *testTxn) Discard(context.Context) error {
	t.useCount++
	if t.useCount == t.failOn {
		return errors.New("DISCARD_ERROR")
	}
	return nil
}

type testDgraphClient struct {
	queryResponse []byte
	alterResponse error
}

func (*testDgraphClient) Login(context.Context, *api.LoginRequest, ...grpc.CallOption) (*api.Response, error) {
	return &api.Response{}, nil
}

func (d *testDgraphClient) Query(context.Context, *api.Request, ...grpc.CallOption) (*api.Response, error) {
	return &api.Response{Json: d.queryResponse}, nil
}

func (d *testDgraphClient) Mutate(context.Context, *api.Mutation, ...grpc.CallOption) (*api.Assigned, error) {
	return &api.Assigned{Uids: map[string]string{"damienstamates": "0x1"}}, nil
}

func (*testDgraphClient) Alter(context.Context, *api.Operation, ...grpc.CallOption) (*api.Payload, error) {
	return &api.Payload{}, nil
}

func (*testDgraphClient) CommitOrAbort(context.Context, *api.TxnContext, ...grpc.CallOption) (*api.TxnContext, error) {
	return &api.TxnContext{}, nil
}

func (*testDgraphClient) CheckVersion(context.Context, *api.Check, ...grpc.CallOption) (*api.Version, error) {
	return &api.Version{}, nil
}

var (
	testValidJSONOutput = []byte(`
	{
		"find": [
			{
				"uid": "0x1"
			}
		]
	}`)
)

// testing data.
var (
	testPersonCorrect = struct {
		Username   string `quirk:"username,unique"`
		Website    string `quirk:"website"`
		AccountAge int    `quirk:"acctage"`
		Email      string `quirk:"email,unique"`
	}{
		Username:   "damienstamates",
		Website:    "github.com",
		AccountAge: 197,
		Email:      "damienstamates@gmail.com",
	}

	testPredValCorrect = predValPairs{
		&predValDat{predicate: "username", value: testPersonCorrect.Username, isUnique: true},
		&predValDat{predicate: "website", value: testPersonCorrect.Website, isUnique: false},
		&predValDat{predicate: "acctage", value: testPersonCorrect.AccountAge, isUnique: false},
		&predValDat{predicate: "email", value: testPersonCorrect.Email, isUnique: true},
	}

	testPersonInvalid = struct {
		Username   string
		Website    string `quirk:"website"`
		AccountAge int    `quirk:"acctage"`
		Email      string `quirk:"email,unique"`
	}{
		Username:   "damienstamates",
		Website:    "github.com",
		AccountAge: 197,
		Email:      "damienstamates@gmail.com",
	}

	testPredValInvalid = predValPairs{
		&predValDat{predicate: "", value: testPersonCorrect.Username, isUnique: false},
		&predValDat{predicate: "website", value: testPersonCorrect.Website, isUnique: false},
		&predValDat{predicate: "acctage", value: testPersonCorrect.AccountAge, isUnique: false},
		&predValDat{predicate: "email", value: testPersonCorrect.Email, isUnique: true},
	}
)
