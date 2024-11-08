package su1

import (
	"context"
	"fmt"
	"os"
)

const filePath = "./su1/tmp/concurrent-write.txt"

func CreateFile(fileName string) *os.File {
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil
	}

	fmt.Println("File created successfully")

	return file
}

func WriteWithOffset(file *os.File, offset int64, chunk []byte) error {
	_, err := file.WriteAt(chunk, offset)

	return err
}

func ConcurrentWrite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file := CreateFile(filePath)
	defer file.Close()

	tasks := make(chan Task)

	done, errors := WorkerPool(ctx, 3, tasks)

	chunk := []byte("hello")

	for i := 0; i < 3; i++ {
		idx := i

		tasks <- func() error {
			err := WriteWithOffset(file, int64(idx*len(chunk)), chunk)

			if err != nil {
				return err
			}

			return nil
		}
	}

	close(tasks)

	select {
	case err := <-errors:
		if err != nil {
			fmt.Printf(err.Error())
			cancel()
		}
	case <-done:
		fmt.Println("completed without errors")
	}
}
