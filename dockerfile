# STAGE 0: Contruct build base
FROM golang:1.16-stretch as builder_base

WORKDIR /order

COPY go.mod .
COPY go.sum .

RUN go mod download

# STAGE 1: Build binaries
FROM builder_base as builder
COPY . /order_service
WORKDIR /order_service

RUN go build -a -installsuffix cgo -o api ./cmd/api

# STAGE 2: Build Workflow Worker
FROM alpine as workflowworker
COPY --from=builder /order_service/api /go/bin/api
RUN apk add -U --no-cache ca-certificates
ENTRYPOINT /go/bin/api
