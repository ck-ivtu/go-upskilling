package su3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type QuickChartClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewQuickChartClient() *QuickChartClient {
	return &QuickChartClient{
		httpClient: InitClient(),
		baseURL:    "https://quickchart.io/chart",
	}
}

func (c *QuickChartClient) GenerateChart(ctx context.Context, chart string, outputFileName string) error {
	chartURL, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}

	query := chartURL.Query()
	query.Set("c", chart)
	chartURL.RawQuery = query.Encode()

	fmt.Printf("Requesting chart from URL: %s\n", chartURL.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, chartURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	file, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("File generated: %s\n", outputFileName)

	return nil
}

func QuickChart() {
	client := NewQuickChartClient()
	ctx := context.Background()

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

	err := client.GenerateChart(ctx, chart, "chart.png")
	if err != nil {
		panic(err)
	}
}
