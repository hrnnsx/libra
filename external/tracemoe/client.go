package tracemoe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const endpoint = "https://api.trace.moe/search"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Search(
	ctx context.Context,
	imageURL string,
) (*SearchResponse, error) {
	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	query := requestURL.Query()
	query.Set("url", imageURL)
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"trace.moe returned status %d",
			resp.StatusCode,
		)
	}

	var result SearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf(
			"trace.moe error: %s",
			result.Error,
		)
	}

	return &result, nil
}
