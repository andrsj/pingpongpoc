package resty

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"pingpongpoc/internal/constants"
)

type Client struct {
	restyClient *resty.Client
	log         *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	return &Client{
		restyClient: resty.New(),
		log:         logger,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	c.log.Info("Calling ping",
		"Query Params", fmt.Sprintf("%s=%s", constants.QueryParamKey1, constants.QueryParamValue1),
		"url", constants.TCPLocalBaseURL,
	)

	request := c.restyClient.R().
		SetContext(ctx).
		SetQueryParam(constants.QueryParamKey1, constants.QueryParamValue1)

	response, err := request.Get(constants.TCPLocalBaseURL)
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
