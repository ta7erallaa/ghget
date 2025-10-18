// Package client
package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *Client) FetchFile(url string) (io.ReadCloser, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if err := checkStatus(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	return resp.Body, nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("file not found (404) - check user, repo, branch or file name")
	case http.StatusForbidden:
		return fmt.Errorf("access forbidden (403) - check repository visibility")
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized (401) - authentication required")
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limited (429) - try again later")
	default:
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}
}
