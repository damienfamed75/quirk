package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateMultiDupleNode(ctx context.Context, dg DgraphClient,
	dat []*DupleNode, uidMap map[string]string) error {
	// Create waitgroup and channels.
	var (
		limit = maxWorkers
		wg    sync.WaitGroup
		m     sync.Mutex
		quit  = make(chan bool)
		read  = make(chan interface{}, len(dat))
		done  = make(chan error)
	)

	if len(dat) < maxWorkers {
		limit = len(dat)
	}

	// Launch workers.
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go mutationWorker(ctx, dg, &wg, &m, c.mutateSingleDupleNode, c.logger,
			uidMap, read, quit, done)
	}

	// Send data to workers via channel.
	for _, d := range dat {
		read <- d
	}

	close(read)

	return launchWorkers(limit, &wg, done, quit)
}
