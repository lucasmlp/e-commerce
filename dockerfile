# STAGE 0: Contruct build base
FROM golang:1.16-stretch as builder_base

WORKDIR /order

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod .
COPY go.sum .

RUN go mod download

# STAGE 1: Build binaries
FROM builder_base as builder
COPY . /order_service
WORKDIR /order_service

RUN go build -a -installsuffix cgo -o workflowworker ./cmd/workflowworker

# STAGE 2: Build Workflow Worker
FROM alpine as workflowworker
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
COPY --from=builder /order_service/workflowworker /go/bin/workflowworker
RUN apk add -U --no-cache ca-certificates
ENTRYPOINT /go/bin/workflowworker
