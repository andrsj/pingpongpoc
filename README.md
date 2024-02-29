# pingpongpoc

## Testing

### HTTP TCP-IP

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

### HTTP Unix

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

### Unix socket

```shell
printf "ping" | socat -t10 UNIX:/tmp/pingpong2.sock -
Response: pong
```
