FROM golang:1.25.4-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/go-rest-api/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/config/docker.yaml ./docker.yaml
COPY --from=builder /app/docs/openapi.yml ./docs/openapi.yml

EXPOSE 8080
