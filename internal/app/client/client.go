package client

import (
	"context"
	"errors"
	"fmt"

	"goapp/internal/pkg/clientsession"

	"golang.org/x/sync/errgroup"
)

// Start opens concurrent WebSocket sessions and logs every received counter value.
func Start(ctx context.Context, connections int) error {
	// Validate the input (url, connections). Is it really needed here?

	group := &errgroup.Group{}

	// Capture errors separately so that we can log them later.
	errs := []error{}

	for id := 0; id < connections; id++ {
		session := clientsession.New(id)

		// Start one goroutine per connection.
		group.Go(func() error {
			err := session.Start(ctx)
			if errors.Is(err, context.Canceled) {
				return nil
			}

			errs = append(errs, err)

			return err
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("at least one connection failed: %w", errors.Join(errs...))
	}

	return nil
}
