#---Build stage---
FROM golang AS builder
COPY . /go/src/
WORKDIR /go/src/cmd/service-pdf-compose-server

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s' -o /go/bin/service

#---Final stage---
FROM alpine:latest
COPY --from=builder /go/bin/service /go/bin/service
CMD /go/bin/service --port 8090 --host '0.0.0.0'