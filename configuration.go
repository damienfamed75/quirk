package quirk

import "github.com/damienfamed75/yalp"

// ClientConfiguration is used to pass in options
// to change the client and customize it to the user's liking.
type ClientConfiguration func(*Client)

// WithLogger sets the logger used by the quirk client.
// By default this is quirk.NewNilLogger.
func WithLogger(l yalp.Logger) ClientConfiguration {
	return func(c *Client) {
		c.logger = l
	}
}

// WithPredicateKey sets the field(predicate) that will
// be used to label inserted nodes. By default this is "name"
func WithPredicateKey(predicateName string) ClientConfiguration {
	return func(c *Client) {
		c.predicateKey = predicateName
	}
}

// WithTemplate sets the field in the Quirk client that
// uses a progress bar to show the nodes being inserted with multi
// node sets.
func WithTemplate(tmpl string) ClientConfiguration {
	return func(c *Client) {
		c.template = tmpl
	}
}

// WithMaxWorkerCount will set the maximum workers that will be spun when
// using a Multi operation.
func WithMaxWorkerCount(count int) ClientConfiguration {
	return func(c *Client) {
		c.maxWorkerCount = count
	}
}

// WithDebugMode sets the client to return verbose errors when possible.
// For instance, when there is a query error then with debug mode on, the client
// will return an error that contains the entire attempted query with the error.
func WithDebugMode() ClientConfiguration {
	return func(c *Client) {
		c.verbose = true
	}
}
