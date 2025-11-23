FROM golang:1.25.4-trixie AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ./migrations /migrations
