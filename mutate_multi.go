package quirk

import (
	"context"
	"sync"

	"github.com/cheggaaa/pb/v3"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/y"
)

// mutateMulti is used for all kinds of mutating any multiple type.
func (c *Client) mutateMulti(ctx context.Context, dg *dgo.Dgraph,
	dat []interface{}, uidMap map[string]UID, mutateFunc mutateSingle) error {
	// Create waitgroup and channels.
	var (
		// wg     sync.WaitGroup
		m      sync.Mutex
		limit  = c.maxWorkerCount
		datLen = len(dat)
		quit   = make(chan bool)
		read   = make(chan interface{}, datLen)
		done   = make(chan error)
	)

	// If there is less data than the max worker count.
	if len(dat) < maxWorkers {
		limit = datLen
	}

	// Create the progress bar.
	bar := pb.ProgressBarTemplate(c.template).Start(datLen)
	bar.SetWidth(bar.Width()/2 + bar.Width()/4)

	// pkg is the more non-focused items that when reading through as a new
	// user, you don't need to focus on as much as some others. For example
	// the user would want to see the path of the UID map and finding it
	// will be easier with the lesser amount of parameters.
	pkg := &workerPackage{
		dg:                 dg,
		m:                  &m,
		mutateSingleStruct: mutateFunc,
		logger:             c.logger,
		bar:                bar,
	}

	// Launch workers.
	for i := 0; i < limit; i++ {
		// go mutationWorker(ctx, dg, &m, mutateFunc, c.logger, bar,
		// 	uidMap, read, quit, done)
		go mutationWorker(ctx, pkg, uidMap, read, quit, done)
	}

	// Send data to workers via channel.
	for _, d := range dat {
		read <- d
	}

	close(read)

	return launchWorkers(limit, bar, done, quit)
}

func launchWorkers(limit int, bar *pb.ProgressBar,
	done chan error, quit chan bool) error {

	var err error
	// Wait for workers to finish.
	// receive results from channel.
	for i := 0; i < limit; i++ {
		// Read the write error from the done channel.
		werr := <-done
		if werr != nil {
			err = werr
			// Close the quit channel to stop the rest of the workers.
			close(quit)
			i = limit
		}
	}

	// Wait for all the workers to finish.
	bar.Finish()

	return err
}

func mutationWorker(ctx context.Context, pkg *workerPackage,
	uidMap map[string]UID, read chan interface{}, quit chan bool, done chan error) {
	// Defer that the waitgroup is finished.
	var err error

	// For each signal received in read channel.
ReadLoop:
	for data := range read {
		// Loop through until a definitive error or success message
		// is received from a mutation.
	Forever:
		for {
			// MutateSingleStruct with received struct.
			_, mutErr := pkg.mutateSingleStruct(ctx, pkg.dg, data, uidMap, pkg.m)

			switch mutErr {
			case nil:
				// If the node was successful then continue to next node.
				break Forever
			case y.ErrAborted:
				// If the transaction was aborted then retry.
			default:
				err = mutErr
				break ReadLoop
			}
		}

		// Increment the progress bar once the node is either successfully added
		// or successfully updated.
		pkg.bar.Increment()
	}

	// Mark done.
	select {
	// Insert the err variable into done.
	// Note: err can be nil or an actual error.
	case done <- err:
		return
	// If a signal was given from quit then return immediately.
	case <-quit:
		return
	}
}
