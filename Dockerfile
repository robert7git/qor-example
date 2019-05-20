FROM golang:1.12.5-alpine3.9
ENV PORT 7000
EXPOSE 7000

WORKDIR /go/src/qor-example
COPY . .

RUN go get

ENTRYPOINT go run main.go
