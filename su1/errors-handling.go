package su1

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func ErrorsHandling() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := make(chan Task)

	done, errors := WorkerPool(ctx, 3, tasks)

	for i := 0; i < 10; i++ {
		idx := i

		tasks <- func() error {
			<-time.After(time.Second * 1)

			if idx%2 == 0 {
				return fmt.Errorf("an error occurred within a task %d\n", idx)
			}

			println("Task " + strconv.Itoa(idx) + " is completed")

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
