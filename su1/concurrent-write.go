package su1

import (
	"context"
	"fmt"
	"os"
	"sync"
)

func CreateFile(fileName string) *os.File {
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil
	}

	fmt.Println("File created successfully")

	return file
}

func WriteWithOffset(file *os.File, index int, chunk string, mu *sync.Mutex) {
	offset := int64(index * len(chunk))
	mu.Lock()
	defer mu.Unlock()

	_, err := file.WriteAt([]byte(chunk+"\n"), offset)

	if err != nil {
		panic(err)
	}
}

func ConcurrentWrite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file := CreateFile("./su1/tmp/concurrent-write.txt")
	defer file.Close()

	tasks := make(chan Task)

	done, errors := WorkerPool(ctx, 3, tasks)

	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		idx := i

		tasks <- func() error {
			WriteWithOffset(file, idx, "hello", &mu)

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
