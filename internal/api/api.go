// internal/api/api.go

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/guarzo/zkillanalytics/internal/model"
)

const (
	maxRetries = 5
	baseDelay  = 1 * time.Second
	maxDelay   = 32 * time.Second
)

func RetryWithExponentialBackoff(operation func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	delay := baseDelay

	for i := 0; i < maxRetries; i++ {
		if result, err = operation(); err == nil {
			return result, nil
		}

		var customErr *CustomError
		if !errors.As(err, &customErr) || (customErr.StatusCode != http.StatusServiceUnavailable && customErr.StatusCode != http.StatusGatewayTimeout) && customErr.StatusCode != http.StatusInternalServerError {
			break
		}

		if i == maxRetries-1 {
			break
		}

		jitter := time.Duration(rand.Int63n(int64(delay)))
		time.Sleep(delay + jitter)

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
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
