package su1

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const fileUrl = "https://example.com/"
const downloadFilePath = "./su1/tmp/concurrent-download.html"

func GetFileSize(url string) (int64, error) {
	client := http.Client{}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
}

func DownloadChunk(url string) func(start int64, end int64) ([]byte, error) {
	client := http.Client{}

	return func(start int64, end int64) ([]byte, error) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}
}

func ConcurrentDownload() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := make(chan Task)
	done, errors := WorkerPool(ctx, 3, tasks)

	file := CreateFile(downloadFilePath)
	defer file.Close()

	fileSize, err := GetFileSize(fileUrl)
	if err != nil {
		panic(err)
	}

	println("Filesize: ", fileSize)

	chunkSize := int64(100)
	downloadChunk := DownloadChunk(fileUrl)

	for off := int64(0); off < fileSize; off += chunkSize {
		offset := off

		tasks <- func() error {
			chunk, err := downloadChunk(offset, offset+chunkSize)
			if err != nil {
				return err
			}

			println("offset", offset, "\n\n\n\n", string(chunk), "\n\n\n\n")

			err = WriteWithOffset(file, offset, chunk)
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
