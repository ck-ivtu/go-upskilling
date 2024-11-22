package su3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileAPIClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewFileAPIClient() *FileAPIClient {
	return &FileAPIClient{
		httpClient: InitClient(),
		baseURL:    "http://localhost:8080",
	}
}

func (c *FileAPIClient) UploadFile(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/files", c.baseURL), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	var result map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["file ID"], nil
}

func FileClient() {
	ctx := context.Background()

	chartClient := NewQuickChartClient()
	chart := `{
		"type": "bar",
		"data": {
			"labels": ["Red", "Blue", "Yellow", "Green", "Purple", "Orange"],
			"datasets": [{
				"label": "Votes",
				"data": [12, 19, 3, 5, 2, 3]
			}]
		}
	}`
	tempFile := "chart.png"

	err := chartClient.GenerateChart(ctx, chart, tempFile)
	if err != nil {
		panic(err)
	}

	fileClient := NewFileAPIClient()

	fileID, err := fileClient.UploadFile(ctx, tempFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File uploaded: %s\n", fileID)

	err = os.Remove(tempFile)
	if err != nil {
		panic(err)
	}
}
