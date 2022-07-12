# Go-Translate

Go-Translate implements a translation relay server for Brave Translate. Requests are sent from Brave browser through a proxy network to Go-Translate, and forwarded to the remote Translation Service.

```
Browser/Proxy (Source) <=> Go-Translate <=> Translation Service (Target)
```

Go-Translate relays the following requests:

```
(Source) GET <go-translate>/translate_a/l => GET <target>/get-languages
(Source) POST <go-translate>/translate_a/t => POST <target>/translate
```

Additionally Go-Translate serves static files necessary for page translation:

```
(Source) GET <go-translate>/static/v1/element.js
(Source) GET <go-translate>/static/v1/js/element/main.js
(Source) GET <go-translate>/static/v1/css/translateelement.css
```

## Dependencies

- Install Go 1.12 or later.
- Dependencies are managed by go modules.
- `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`

## Setup

- Clone `git clone git@github.com:brave/go-translate.git`
- Build and run: `make build`
- Run linter: `make lint`
- Run tests: `make test`
