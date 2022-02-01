FROM golang:1.16-alpine


RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o go-translate
EXPOSE 8195

CMD ["/app/go-translate"]
