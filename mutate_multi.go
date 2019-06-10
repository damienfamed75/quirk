package quirk

import (
	"context"
	"sync"
)

func (c *Client) mutateMultiStruct(ctx context.Context, dg DgraphClient,
	dat []interface{}, uidMap map[string]string) error {
	// Create waitgroup and channels.
	var (
		limit = maxWorkers
		wg    sync.WaitGroup
		m     sync.Mutex
		quit  = make(chan bool)
		read  = make(chan interface{}, len(dat))
		write = make(chan map[string]string, len(dat))
		done  = make(chan error)
	)

	if len(dat) < maxWorkers {
		limit = len(dat)
	}

	// Launch workers.
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go mutationWorker(ctx, dg, &wg, &m, c.mutateSingleStruct, uidMap,
			read, quit, done)
	}

	// Send data to workers via channel.
	for _, d := range dat {
		read <- d
	}

	close(read)

	return launchWorkers(limit, &wg, write, done, quit)
}

func launchWorkers(limit int, wg *sync.WaitGroup, write chan map[string]string,
	done chan error, quit chan bool) error {
	var err error

	// Wait for workers to finish.
	// receive results from channel.
	for i := 0; i < limit; i++ {
		select {
		case err = <-done:
			if err != nil {
				close(quit)
				i = limit
			}
		}
	}

	wg.Wait()

	return err
}

func mutationWorker(ctx context.Context, dg DgraphClient, wg *sync.WaitGroup,
	m *sync.Mutex, mutateSingleStruct mutateSingle, uidMap map[string]string,
	read chan interface{}, quit chan bool, done chan error) {
	// Defer that the waitgroup is finished.
	defer wg.Done()
	var err = error(nil)

	// For each signal received in read channel.
	for data := range read {
		// MutateSingleStruct with received struct.
		mutErr := mutateSingleStruct(ctx, dg, data, uidMap, m)
		if mutErr != nil {
			err = mutErr
			break
		}
	}

	// Mark done.
	select {
	case done <- err:
		return
	case <-quit:
		return
	}
}
