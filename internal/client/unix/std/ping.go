package std

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"

	"pingpongpoc/internal/constants"
)

type Client struct {
	httpClient *http.Client

	log *slog.Logger
}

func NewClient(socketPath string, logger *slog.Logger) *Client {
	//nolint:exhaustruct
	transport := &http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			connection, err := net.Dial(constants.NetworkType, socketPath)
			if err != nil {
				return nil, fmt.Errorf("failed dialing unix socket '%s': %w", socketPath, err)
			}

			return connection, nil
		},
	}

	return &Client{
		httpClient: &http.Client{Transport: transport}, //nolint:exhaustruct
		log:        logger,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	c.log.Info("Calling ping",
		"Query Params", fmt.Sprintf("%s=%s", constants.QueryParamKey1, constants.QueryParamValue1),
		"url", constants.UnixLocalBaseURL,
	)

	fullURL, err := c.formatFullURL()
	if err != nil {
		return fmt.Errorf("failed formatting full URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed making GET request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading body error: %w", err)
	}

	c.log.Info("Response",
		"status", response.Status,
		"body", responseBody,
		"proto", response.Proto,
		"headers", response.Header,
	)

	return nil
}

// formatFullURL makes fullURL := url + "?" + queryParamKey + "=" + queryParamValue
func (c Client) formatFullURL() (string, error) {
	u, err := url.Parse(constants.UnixLocalBaseURL) //nolint:varnamelen
	if err != nil {
		return "", fmt.Errorf("parsing base URL error: %w", err)
	}

	query := u.Query()
	query.Set(constants.QueryParamKey1, constants.QueryParamValue1)
	u.RawQuery = query.Encode()

	c.log.Info("Full URL", "url", u.String())

	return u.String(), nil
}
