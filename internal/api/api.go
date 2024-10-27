// internal/api/api.go

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gambtho/zkillanalytics/internal/model"
)

// RetryWithExponentialBackoff retries the provided function with exponential backoff.
func RetryWithExponentialBackoff(operation func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	backoff := time.Second

	for i := 0; i < 5; i++ {
		result, err = operation()
		if err == nil {
			return result, nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, err
}

// GetPageData fetches a single page of data given a URL with context.
func GetPageData(ctx context.Context, client *http.Client, url string) ([]model.KillMail, error) {
	result, err := RetryWithExponentialBackoff(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var killMails []model.KillMail
		if err := json.Unmarshal(body, &killMails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		return killMails, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]model.KillMail), nil
}
