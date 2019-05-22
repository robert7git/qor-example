# step 1: build
FROM golang:1.12.5-alpine3.9 as build-step

# for go mod download
RUN apk add --update --no-cache build-base ca-certificates git

RUN mkdir /go-app
WORKDIR /go-app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/qor-example

# -----------------------------------------------------------------------------
# step 2: exec
FROM scratch

COPY --from=build-step /go/bin/qor-example /go/bin/qor-example
COPY app .
COPY public .
COPY config .

CMD ["/go/bin/qor-example"]

