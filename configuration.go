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