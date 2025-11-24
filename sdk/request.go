package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

func (c *Client) do(ctx context.Context, method, url string, body []byte) (*utils.Response, error) {
	var lastErr error
	var lastResp *utils.Response

	for i := 0; i <= c.config.MaxRetries; i++ {
		if i > 0 {
			delay := time.Duration(math.Pow(2, float64(i-1))*100) * time.Millisecond
			// Respect context cancellation during backoff
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", utils.ErrInvalidRequest, err)
		}

		if c.config.Token == "" {
			return nil, utils.ErrInvalidConfiguration
		}

		req.Header.Set("Authorization", "Bearer "+c.config.Token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Read body and close immediately (not with defer in loop!)
		respBody, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close() // Always close, even if ReadAll fails (error ignored - we already have the body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 300 {
			lastResp = &utils.Response{
				Status:   utils.StatusError,
				Error:    string(respBody),
				HTTPCode: resp.StatusCode,
			}

			if resp.StatusCode == http.StatusUnauthorized {
				if c.isKeycloakAuthMethodConfigured() {
					if _, err := c.refreshToken(ctx); err == nil {
						continue // Retry with the new token
					}
				}
				return lastResp, utils.ErrAuthenticationFailed
			}

			if resp.StatusCode == http.StatusForbidden {
				return lastResp, utils.ErrPermissionDenied
			}

			if resp.StatusCode == http.StatusNotFound {
				return lastResp, utils.ErrNotFound
			}

			// Do not retry on other 4xx client errors
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return lastResp, fmt.Errorf("%w: %s", utils.ErrInvalidRequest, string(respBody))
			}

			lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
			continue
		}

		var parsedBody any
		if err := json.Unmarshal(respBody, &parsedBody); err != nil {
			lastErr = fmt.Errorf("failed to parse response body: %w", err)
			continue
		}

		return &utils.Response{
			Status:   utils.StatusOK,
			Data:     parsedBody,
			HTTPCode: resp.StatusCode,
		}, nil
	}

	if lastResp != nil {
		return lastResp, fmt.Errorf("max retries exceeded, last response was: %s", lastResp.Error)
	}

	return nil, fmt.Errorf("max retries exceeded, last error: %w", lastErr)
}
