package resty

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"

	"pingpongpoc/internal/constants"
)

type Client struct {
	restyClient *resty.Client

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

	client := resty.New()

	client.SetTransport(transport)
	client.SetScheme(constants.SchemeType)

	return &Client{
		restyClient: client,
		log:         logger,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	c.log.Info("Calling ping",
		"Query Params", fmt.Sprintf("%s=%s", constants.QueryParamKey1, constants.QueryParamValue1),
		"url", constants.UnixLocalBaseURL,
	)

	request := c.restyClient.R().
		SetContext(ctx).
		SetQueryParam(constants.QueryParamKey1, constants.QueryParamValue1)

	response, err := request.Get(constants.UnixLocalBaseURL)
	if err != nil {
		return fmt.Errorf("failed making GET request: %w", err)
	}

	c.log.Info("Response",
		"status", response.Status(),
		"body", response.Body(),
		"proto", response.Proto(),
		"time", response.Time(),
		"received at", response.ReceivedAt(),
		"header", response.Header(),
	)

	requestTraceInfo := response.Request.TraceInfo()

	// TODO check the real trace info
	// at the moment of using localhost
	// there is no needs to check this info
	c.log.Info("Request Trace Info",
		"DNS lookup", requestTraceInfo.DNSLookup,
		"connection time", requestTraceInfo.ConnTime,
		"TCP connection time", requestTraceInfo.TCPConnTime,
		"TLS handshake", requestTraceInfo.TLSHandshake,
		"server time", requestTraceInfo.ServerTime,
		"response time", requestTraceInfo.ResponseTime,
		"total time", requestTraceInfo.TotalTime,
		"is connection reused", requestTraceInfo.IsConnReused,
		"ss connection was Idle", requestTraceInfo.IsConnWasIdle,
		"connection Idle time", requestTraceInfo.ConnIdleTime,
		"request attempt", requestTraceInfo.RequestAttempt,
		// Be careful, panic if this is local XD
		// "remote address", requestTraceInfo.RemoteAddr.String(),
	)

	return nil
}
