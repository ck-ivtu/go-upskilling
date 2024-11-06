package su1

import (
	"context"
	"strconv"
	"sync"
)

func DataRace() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mu := sync.Mutex{}

	tasks := make(chan Task)

	done, errors := WorkerPool(ctx, 3, tasks)

	res := make(map[string]int)

	for i := 0; i < 3; i++ {
		idx := i

		tasks <- func() error {
			mu.Lock()
			res[strconv.Itoa(idx)] = idx
			mu.Unlock()

			return nil
		}
	}

	close(tasks)

	for err := range errors {
		if err != nil {
			cancel()
			break
		}
	}

	<-done
}
