package constants

import "time"

const (
	TimeoutForPinging = 10 * time.Second

	PathToSocketHTTPUnixServer = "/tmp/pingpong1.sock"
	PathToSocketUnixServer     = "tmp/pingpong2.sock"

	ShutdownTimeout = 5 * time.Second
	ReadHeadTimeout = 5 * time.Second

	HTTPListenAddressTCPServer = ":8080"

	QueryParamKey1   = "input"
	QueryParamValue1 = "ping"

	TCPLocalBaseURL  = "http://localhost:8080/"
	UnixLocalBaseURL = "http://localhost/"

	SchemeType  = "http"
	NetworkType = "unix"
)
