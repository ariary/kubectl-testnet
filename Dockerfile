# syntax=docker/dockerfile:1
FROM golang:1.20-alpine3.17 AS builder

WORKDIR /build
COPY . .
RUN rm -rf manifests

RUN make before.build
RUN make build.testnet

FROM alpine:3.18.5

RUN apk --no-cache update && \
    apk --no-cache upgrade && \
    apk --no-cache add ca-certificates make vim curl tcpdump bash

RUN useradd -ms /bin/bash app
COPY --chown=app --from=builder /build/testnet /app/testnet
COPY --chown=app --from=builder /build/Makefile /app/Makefile

WORKDIR /app
USER app
