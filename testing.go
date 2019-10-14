package quirk

import (
	"context"
	"errors"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

var (
	// Ensure testDgraphClient fits the api interface.
	// This is for any direct compilation errors to make the pipeline fail faster.
	// Very optimistic of me :)
	_ api.DgraphClient = &testDgraphClient{}
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

// Reset is empty because testBuilder doesn't hold any information in it in the first place.
func (*testBuilder) Reset() {
	// This function is just for testing.
}

type testDgraphClient struct {
	queryUseCount int
	failQueryOn   int
	queryResponse []byte
	alterResponse error
	shouldAbort   bool
}

func (*testDgraphClient) Login(context.Context, *api.LoginRequest, ...grpc.CallOption) (*api.Response, error) {
	return &api.Response{}, nil
}

func (d *testDgraphClient) Query(context.Context, *api.Request, ...grpc.CallOption) (*api.Response, error) {
	d.queryUseCount++
	if d.queryUseCount == d.failQueryOn {
		return &api.Response{}, errors.New("QUERY_ERROR")
	}
	// Changed Response to return a UID Map like the Mutate function.
	// This is because of the latest major change to dgo and how transactions
	// interact with the api client when mutating nodes.
	// See Issue #17 for more info.
	return &api.Response{Json: d.queryResponse,
		Uids: map[string]string{"damienstamates": "0x1"}}, nil
}

// api.Response isn't api.Assigned because of major version changes to dgo
// See Issue #17 for more info.
func (d *testDgraphClient) Mutate(context.Context, *api.Mutation, ...grpc.CallOption) (*api.Response, error) {
	if d.shouldAbort {
		return &api.Response{}, dgo.ErrAborted
	}
	return &api.Response{Uids: map[string]string{"damienstamates": "0x1"}}, nil
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

	testPredValCorrect = &DupleNode{
		Duples: []Duple{
			{Predicate: "username", Object: testPersonCorrect.Username, IsUnique: true},
			{Predicate: "website", Object: testPersonCorrect.Website, IsUnique: false},
			{Predicate: "acctage", Object: testPersonCorrect.AccountAge, IsUnique: false},
			{Predicate: "email", Object: testPersonCorrect.Email, IsUnique: true},
		},
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

	testPredValInvalid = &DupleNode{
		Duples: []Duple{
			{Predicate: "", Object: testPersonCorrect.Username, IsUnique: false},
			{Predicate: "website", Object: testPersonCorrect.Website, IsUnique: false},
			{Predicate: "acctage", Object: testPersonCorrect.AccountAge, IsUnique: false},
			{Predicate: "email", Object: testPersonCorrect.Email, IsUnique: true},
		},
	}
)
