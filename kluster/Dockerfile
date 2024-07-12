# syntax=docker/dockerfile:1.7.0

ARG BUILDPLATFORM="linux/arm64"

FROM --platform=$BUILDPLATFORM golang:1.22 AS builder

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN  go mod download

COPY . .
RUN  go build -o monokube ./cmd/server


FROM alpine:3.20.1 AS release
LABEL co.ultary.image.authors="ghilbut@gmail.com"

EXPOSE 9090
EXPOSE 50051

CMD ["monokube", "serve"]

COPY --from=builder /app/monokube /usr/local/bin/monokube
