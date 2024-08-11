# Clients Logistics Engine API

**What is Our  Example Project**

This is example of repository that hosts the source code for a client(s) designed to generate a certain number of requests (N), as specified by its operations, encoded using protocol buffers.  It is intentionally divided client and server source-code between two different repositories, to give a sense of interaction two real distributed systems.
You can find the complete server example [here](https://github.com/ivanbulyk/logistics_engine_api).

So to see everything in action, you can test it by running our server first:

```text
$ go run ./cmd/logistics/ main.go
```
or just

```text
$ make 
```

Then in other terminal, we run our client, in project root:

```text
$ go run ./cmd/logistics/ main.go
```
or just

```text
$ make 
```

The server should wait infinitely, emitting logs on calls, and the client should be returning without any error on the terminal. Then you want to hit the localhost:50051 LogisticsEngineAPI/MetricsReport, with any gRPC client to see the calculations result.