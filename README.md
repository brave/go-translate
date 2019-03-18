# Translation relay server for brave

`go-translate` implements a translation relay server for use in brave-core written in Go.

The intended audience for this server is all users of brave-core.

The translation relay server supports 2 endpoints

1) The `POST /translate` endpoint processes translate requests in Google format, sends corresponding requests in Microsoft format to Microsoft translate server, then returns responses in Google format back to the brave-core client.

2) The `GET /language` endpoint processes requests of getting the support language list in Google format, sends corresponding requests in Microsoft format to Microsoft translate server, then returns responses in Google format back to the brave-core client.

There are also a few static resources requested during in-page translation will be handled by go-translate and will be proxied through a brave to avoid introducing direct connection to any Google server.


## Dependencies

- Install Go 1.10 or later.
- `dep` is used to install the Go dependencies.
- `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`

## Setup

```
go get -d github.com/brave/go-translate
cd ~/go/src/github.com/brave/go-translate
dep ensure
```

## Run lint:

`make lint`

## Run tests:

`make test`

## Build and run go-translate:

`make build`
