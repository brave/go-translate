FROM golang:1.18-alpine as builder

# put certs in builder image
RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
RUN apk add make
RUN apk add build-base
RUN apk add git
RUN apk add bash

ARG VERSION
ARG BUILD_TIME
ARG COMMIT

WORKDIR /src
COPY . ./

RUN chown -R nobody:nobody /src/
RUN mkdir /.cache
RUN chown -R nobody:nobody /.cache

USER nobody
RUN go mod download

RUN make build

FROM alpine:3.15 as base

# put certs in artifact from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/go-translate /bin/

EXPOSE 8195
CMD ["/bin/go-translate"]
