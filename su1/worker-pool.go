package su1

import (
	"context"
	"fmt"
	"sync"
)

type Task func() error

func WorkerPool(ctx context.Context, maxWorkers int, tasks <-chan Task) (<-chan struct{}, <-chan error) {
	var wg sync.WaitGroup

	done := make(chan struct{})
	errors := make(chan error, maxWorkers)

	worker := func(id int) {
		fmt.Printf("Starting instance %d\n", id)

		defer (func() {
			fmt.Printf("Stopping instance %d\n", id)

			wg.Done()
		})()

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d stopped due to context cancellation\n", id)
				return
			case task, ok := <-tasks:
				if !ok {
					fmt.Printf("Worker %d stopped due to closed channel\n", id)
					return
				}

				err := task()

				if err != nil {
					select {
					case errors <- err:
					default:
					}
				}
			}
		}
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go worker(i + 1)
	}

	go func() {
		wg.Wait()

		close(done)
		close(errors)
	}()

	return done, errors
}
