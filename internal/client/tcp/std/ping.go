package std

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"pingpongpoc/internal/constants"
)

type Client struct {
	log *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	return &Client{
		log: logger,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	c.log.Info("Calling ping",
		"Query Params", fmt.Sprintf("%s=%s", constants.QueryParamKey1, constants.QueryParamValue1),
		"url", constants.TCPLocalBaseURL,
	)

	fullURL, err := c.formatFullURL()
	if err != nil {
		return fmt.Errorf("failed formatting full URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Read and log the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	c.log.Info("Response",
		"status", resp.Status,
		"headers", resp.Header,
		"body", string(body),
	)

	return nil
}

// formatFullURL makes fullURL := url + "?" + queryParamKey + "=" + queryParamValue
func (c Client) formatFullURL() (string, error) {
	u, err := url.Parse(constants.TCPLocalBaseURL) //nolint:varnamelen
	if err != nil {
		return "", fmt.Errorf("parsing base URL error: %w", err)
	}

	query := u.Query()
	query.Set(constants.QueryParamKey1, constants.QueryParamValue1)
	u.RawQuery = query.Encode()

	c.log.Info("Full URL", "url", u.String())

	return u.String(), nil
}
