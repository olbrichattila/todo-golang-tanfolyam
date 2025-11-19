# Stage 1: Builder
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod download
COPY internal/ ./internal

RUN go install github.com/olbrichattila/godbmigrator_cmd/cmd/migrator@latest
RUN go build -o todo .

# Stage 2: Runtime
FROM debian:bookworm-slim


WORKDIR /app
COPY --from=builder /app/todo /app/todo
COPY --from=builder /go/bin/migrator /app/migrator

COPY migrations/ /app/migrations/
COPY static/ /app/static/
COPY templates/ /app/templates/

COPY .env.prod /app/.env
COPY .env.migrator.prod /app/.env.migrator

COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh \
    && mkdir /app/session \
    && mkdir /app/data


CMD ["/app/start.sh"]
