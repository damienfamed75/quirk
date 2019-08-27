package quirk

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/damienfamed75/yalp"

	"github.com/cheggaaa/pb/v3"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/y"
)

var (
	lastStatus   = time.Now()
	successCount uint64
	retryCount   uint64
)

// mutateMulti is used for all kinds of mutating any multiple type.
func (c *Client) mutateMulti(ctx context.Context, dg *dgo.Dgraph,
	dat []interface{}, uidMap map[string]UID, mutateFunc mutateSingle) error {
	// Create waitgroup and channels.
	var (
		limit = c.maxWorkerCount
		wg    sync.WaitGroup
		m     sync.Mutex
		quit  = make(chan bool)
		read  = make(chan interface{}, len(dat))
		done  = make(chan error)
	)

	if len(dat) < maxWorkers {
		limit = len(dat)
	}

	bar := pb.ProgressBarTemplate(c.template).Start(len(dat))
	bar.SetWidth(bar.Width()/2 + bar.Width()/4)

	// Launch workers.
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go mutationWorker(ctx, dg, &wg, &m, mutateFunc, c.logger, bar,
			uidMap, read, quit, done)
	}

	// Send data to workers via channel.
	for _, d := range dat {
		read <- d
	}

	close(read)

	return launchWorkers(limit, &wg, bar, done, quit)
}

func launchWorkers(limit int, wg *sync.WaitGroup, bar *pb.ProgressBar,
	done chan error, quit chan bool) error {

	var err error
	// Wait for workers to finish.
	// receive results from channel.
	for i := 0; i < limit; i++ {
		werr := <-done
		if werr != nil {
			err = werr
			close(quit)
			i = limit
		}
	}

	wg.Wait()
	bar.Finish()

	return err
}

func mutationWorker(ctx context.Context, dg *dgo.Dgraph, wg *sync.WaitGroup,
	m *sync.Mutex, mutateSingleStruct mutateSingle, logger yalp.Logger, bar *pb.ProgressBar,
	uidMap map[string]UID, read chan interface{}, quit chan bool, done chan error) {
	// Defer that the waitgroup is finished.
	defer wg.Done()
	var err error

	// For each signal received in read channel.
ReadLoop:
	for data := range read {
		// Loop through until a definitive error or success message
		// is received from a mutation.
	Forever:
		for {
			// MutateSingleStruct with received struct.
			new, mutErr := mutateSingleStruct(ctx, dg, data, uidMap, m)

			switch mutErr {
			case nil:
				if new {
					// If a successful new node was added then count up.
					atomic.AddUint64(&successCount, 1)
				}
				break Forever
			case y.ErrAborted:
				// If the transaction was aborted then retry.
				atomic.AddUint64(&retryCount, 1)
			default:
				err = mutErr
				break ReadLoop
			}
		}

		// Increment the progress bar once the node is either successfully added
		// or successfully updated.
		bar.Increment()
	}

	// Mark done.
	select {
	case done <- err:
		return
	case <-quit:
		return
	}
}
