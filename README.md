# Overview
This project is an example of how to build a rudementary pub-sub REST API. 
Clients can subscribe to receive new information by invoking the API

> GET /subscribe

Many subscribers can connect and each subscriber will receive any new 
published messages.

Subscribers can connect as a web socket or normal HTTP connection. The choice 
is based on the subscriber. Web sockets are bidirectional, but in this case, 
data is pushed to the clients only. When connecting via straight HTTP, data 
will be streamed over the existing socket so clients do not need to rebind.

The publisher can send new messages to the subscribers by invoking the REST API

> POST /publish

Where the body of the POST method is the message to be sent to all subscribers.

# Unit tests
Unit tests were created for the code that managers the pub-sub logic. To run
the unit tests, run `go test ./...`.

# Running
## Local 
* Run `go run main.go` to start the server

## Docker
* Run `docker build -t web-pub-sub:latest .` to build the container
* Run `docker run -it -p 8000:8000 web-pub-sub:latest` to run the server

## HTTP subscriber
* Run `curl -v http://localhost:8000/subscribe` to register subscribers in as many terminals as desired

## Websocket subscriber
Run the following curl command to upgrade the connection to a web socket.
```
curl                                                    \
    --http1.1                                           \
    --include                                           \
    --no-buffer                                         \
    --header "Connection: Upgrade"                      \
    --header "Upgrade: websocket"                       \
    --header "Host: localhost:8000"                     \
    --header "Origin: http://localhost:8000"            \
    --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ=="  \
    --header "Sec-WebSocket-Version: 13"                \
    http://localhost:8000/subscribe
```

## Publisher
* Run `dd if=/dev/random bs=16 count=1 2>/dev/null | base64 | curl -v --data @- http://localhost:8000/publish` to generate random data to be published
