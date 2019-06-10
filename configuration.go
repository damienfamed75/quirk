package quirk

import "github.com/damienfamed75/quirk/logging"

type ClientConfiguration func(*Client)

func WithLogger(l logging.Logger) ClientConfiguration {
	return func(c *Client) {
		c.logger = l
	}
}

func WithPredicateKey(predicateName string) ClientConfiguration {
	return func(c *Client) {
		c.predicateKey = predicateName
	}
}

func WithQuirkName(name string) ClientConfiguration {
	return func(c *Client) {
		c.quirkName = name
	}
}

func WithQuirkRel(relName string) ClientConfiguration {
	return func(c *Client) {
		c.quirkRel = relName
	}
}

func UseUIDIncrementer() ClientConfiguration {
	return func(c *Client) {
		c.insertMode = incrementor
	}
}

func UseReverseEdges() ClientConfiguration {
	return func(c *Client) {
		c.reverseEdge = true
	}
}

func WithStartUID(i uint64) ClientConfiguration {
	setStartUID(i)
	return func(*Client) {}
}
