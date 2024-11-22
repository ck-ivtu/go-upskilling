package su3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type URLShortenerClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewURLShortenerClient() *URLShortenerClient {
	return &URLShortenerClient{
		httpClient: InitClient(),
		baseURL:    "https://cleanuri.com/api/v1/shorten",
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ResultURL string `json:"result_url"`
}

func (c *URLShortenerClient) Shorten(ctx context.Context, longURL string) (string, error) {
	requestBody, err := json.Marshal(ShortenRequest{URL: longURL})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	var shortenResp ShortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&shortenResp); err != nil {
		return "", err
	}

	return shortenResp.ResultURL, nil
}

func ShortenURL() {
	client := NewURLShortenerClient()
	ctx := context.Background()

	longURL := "https://example.com/very/long/url"
	shortURL, err := client.Shorten(ctx, longURL)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Short URL: %s\n", shortURL)
}
