FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /ingest .

FROM alpine:latest

COPY --from=builder /ingest /ingest

EXPOSE 4317

ENTRYPOINT ["/ingest"]
