# Translation server for brave

`go-translate` implements a translation server for use in brave-core written in Go.
It works on top of brave-hosted [Lingvanex On-premise Translation Server](https://lingvanex.com/translationserver/). It gets the requests from browsers in Chromium format,
rewrites them to Lingvanex format, processes the requests and returns the result back in Chromium format.

The audience for this server is all desktop/android brave users.

The translation server supports 2 endpoints

1) The `POST /translate_a/t` endpoint processes translate requests in Chromium format, sends corresponding requests to Lingvanex docker container, then returns responses in Chromium format back to the brave-core client.

2) The `GET /translate_a/l` returns the languages supported by Lingvanex in Chromium format.

go-translate also hosts a few static resources needed for in-page translation.

## Dependencies

- Install Go 1.12 or later.
- Dependencies are managed by go modules.
- `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`

## Setup

```
git clone git@github.com:brave/go-translate.git
cd ~/path-to/go-translate
make build
```

## Run lint:

`make lint`

## Run tests:

`make test`

## Build and run go-translate:

`make build`

## Local debugging

- Generate TLS credentials:
`openssl genrsa -out server.key 2048`
`openssl ecparam -genkey -name secp384r1 -out server.key`
`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`

- Replace `err := srv.ListenAndServe()` with `err := srv.ListenAndServeTLS("server.crt", "server.key")` in `server/server.go`;

- Set LNX_HOST (company VPN should be enabled):
  `export LNX_HOST=http://translate-lnx-dev-a4b82554457afe1c.elb.us-west-2.amazonaws.com:8080/api`

- Launch the local server: `make build`;

- Launch the browser with switches:
`--ignore-certificate-errors --translate-security-origin=https://127.0.0.1:8195/ --translate-script-url=https://127.0.0.1:8195/static/v1/element.js`.

- Other switches can be added if necessary (for example `--enable-features=UseBraveTranslateGo:update-languages/true` );

- Disable Shield on the tested sites (or globally), it cuts requests to localhost. `405 OPTIONS` errors also can be ignored;

- If you have troubles check if you can reach https://127.0.0.1:8195/static/v1/element.js in the tested browser.
