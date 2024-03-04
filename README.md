# Ping Pong POC (Proof of Concept)

Ping Pong POC (Proof of Concept) is a project which include my view on Clean Architecture of simple application
that has different interfaces for simple functionality. This repo has only one purpose - showing the app can have different types
of communication (TCP, Unix, etc). Aiming to keep the business logic independent of external frameworks and tools.
By separating concerns into distinct layers, such as domain logic, use cases, and infrastructure,
it ensures flexibility, maintainability, and testability of the codebase.

## Features

* **TCP Client**: Supports standard Go `net` package and third-party (`resty`) implementations for TCP communications.
* **Unix Socket Client**: Provides implementations for Unix socket communications with both standard Go `net` package and third-party (`resty`) support
* **Servers**: Includes a basic servers setup, supporting both TCP and Unix socket connections.
* **Extensible CLI**: *NOT FINISHED YET! (fuck Cobra, I'm tired to debug this shit)*

## TODO

* Finish Cobra configuration
* File locker (to avoid running several instances of server)

## Project structure

The project is organized into several key directories:

* `cmd/`: Houses the main entry point for the application
* `internal/`: Includes internal packages for the CLI commands, client implementations, domain logic, and more.
  * `cli/`: Core CLI command implementations using Cobra.
  * `client/`: Client logic for TCP and Unix socket communications.
  * `constants/`: Shared constants used across the project.
  * `domain/`: Domain-specific logic, including any business rules. (can be distributed in `internal/`)
  * `errors/`: Custom error definitions.
  * `logging/`: Logging utilities.
  * `service/`: Service layer for business logic. (in other words `usecases`)
  * `transport/`: HTTP and Unix socket server implementations. (can be more, gRPC, GraphQL, Pub/Subs, etc)

## Testing

### Running

#### Server

```shell
go run ./cmd/server/main.go

{"time":"2024-03-04T08:49:08.692463613+02:00","level":"INFO","msg":"Starting server","address":":8080"}
{"time":"2024-03-04T08:49:08.692525322+02:00","level":"ERROR","msg":"Server failed to start","error":"listen unix tmp/pingpong2.sock: bind: no such file or directory"}
{"time":"2024-03-04T08:49:08.692559098+02:00","level":"ERROR","msg":"Error starting Unix server2","error":"failed to listen on socket: listen unix tmp/pingpong2.sock: bind: no such file or directory"}
{"time":"2024-03-04T08:49:08.692715087+02:00","level":"INFO","msg":"Starting server","socketPath":"/tmp/pingpong1.sock"}

. . .
```

#### HTTP TCP Clients

##### 3rd party client (`resty`)

```shell
go run ./cmd/client/tcp/3rd/main.go

time=2024-03-04T08:50:15.013+02:00 level=INFO msg="Calling ping" "Query Params"="input=ping" url=http://localhost:8080/
time=2024-03-04T08:50:21.017+02:00 level=INFO msg=Response status="200 OK" body="{\"message\":\"Response: pong\"}" proto=HTTP/1.1 time=6.00330752s "received at"=2024-03-04T08:50:21.017+02:00 header="map[Content-Length:[28] Content-Type:[application/json] Date:[Mon, 04 Mar 2024 06:50:21 GMT]]"
time=2024-03-04T08:50:21.017+02:00 level=INFO msg="Request Trace Info" "DNS lookup"=0s "connection time"=0s "TCP connection time"=0s "TLS handshake"=0s "server time"=0s "response time"=0s "total time"=0s "is connection reused"=false "ss connection was Idle"=false "connection Idle time"=0s "request attempt"=0
```

##### STD lib

```shell
go run ./cmd/client/tcp/std/main.go

time=2024-03-04T08:51:17.117+02:00 level=INFO msg="Calling ping" "Query Params"="input=ping" url=http://localhost:8080/
time=2024-03-04T08:51:17.117+02:00 level=INFO msg="Full URL" url="http://localhost:8080/?input=ping"
time=2024-03-04T08:51:23.122+02:00 level=INFO msg=Response status="200 OK" headers="map[Content-Length:[28] Content-Type:[application/json] Date:[Mon, 04 Mar 2024 06:51:23 GMT]]" body="{\"message\":\"Response: pong\"}"
```

#### HTTP Unix Clients

##### 3rd party client (`resty`)

```shell
go run ./cmd/client/unix/3rd/main.go

time=2024-03-04T08:52:41.811+02:00 level=INFO msg="Calling ping" "Query Params"="input=ping" url=http://localhost/
time=2024-03-04T08:52:47.813+02:00 level=INFO msg=Response status="200 OK" body="{\"message\":\"Response: pong\"}" proto=HTTP/1.1 time=6.002183334s "received at"=2024-03-04T08:52:47.813+02:00 header="map[Content-Length:[28] Content-Type:[application/json] Date:[Mon, 04 Mar 2024 06:52:47 GMT]]"
time=2024-03-04T08:52:47.814+02:00 level=INFO msg="Request Trace Info" "DNS lookup"=0s "connection time"=0s "TCP connection time"=0s "TLS handshake"=0s "server time"=0s "response time"=0s "total time"=0s "is connection reused"=false "ss connection was Idle"=false "connection Idle time"=0s "request attempt"=0
```

##### STD lib

```shell
go run ./cmd/client/unix/std/main.go

time=2024-03-04T08:53:31.812+02:00 level=INFO msg="Calling ping" "Query Params"="input=ping" url=http://localhost/
time=2024-03-04T08:53:31.812+02:00 level=INFO msg="Full URL" url="http://localhost/?input=ping"
time=2024-03-04T08:53:37.813+02:00 level=INFO msg=Response status="200 OK" body="{\"message\":\"Response: pong\"}" proto=HTTP/1.1 headers="map[Content-Length:[28] Content-Type:[application/json] Date:[Mon, 04 Mar 2024 06:53:37 GMT]]"
```

### By command tools

#### HTTP TCP-IP

```shell
http GET "http://localhost:8080/?input=ping" -vv
GET /?input=ping HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/3.2.1



HTTP/1.1 200 OK
Content-Length: 28
Content-Type: application/json
Date: Thu, 29 Feb 2024 08:49:31 GMT

{
    "message": "Response: pong"
}


Elapsed time: 6.006945335s
```

---

```shell
curl -v "http://localhost:8080/?input=ping"
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /?input=ping HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.6.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 29 Feb 2024 08:42:25 GMT
< Content-Length: 28
< 
* Connection #0 to host localhost left intact
{"message":"Response: pong"}
```

---

```shell
http GET "http://localhost:8080/?input=peng"

HTTP/1.1 400 Bad Request
Content-Length: 93
Content-Type: application/json
Date: Thu, 29 Feb 2024 08:06:18 GMT

{
    "error": "Are you fucking dumb? It's a 'ping-pong' server. Send me just a one word - 'ping'"
}
```

#### HTTP Unix

```shell
curl --unix-socket /tmp/pingpong1.sock "http://localhost/?input=ping"
{"message":"Response: pong"}%
```

---

```shell
curl -v --unix-socket /tmp/pingpong1.sock "http://localhost/?input=ping"
*   Trying /tmp/pingpong1.sock:0...
* Connected to localhost (/tmp/pingpong1.sock) port 80
> GET /?input=ping HTTP/1.1
> Host: localhost
> User-Agent: curl/8.6.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 29 Feb 2024 08:44:27 GMT
< Content-Length: 28
< 
* Connection #0 to host localhost left intact
{"message":"Response: pong"}% 
```

#### Unix socket

```shell
printf "ping" | socat -t10 UNIX:/tmp/pingpong2.sock -
Response: pong
```
