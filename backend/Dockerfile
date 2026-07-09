FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd
RUN go install github.com/jackc/tern/v2@2.4.0

FROM alpine:3.20 AS migrate

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /go/bin/tern /usr/local/bin/tern
COPY internal/db/migrations ./internal/db/migrations
COPY scripts/migrate-entrypoint.sh ./migrate-entrypoint.sh
RUN chmod +x ./migrate-entrypoint.sh

ENTRYPOINT ["./migrate-entrypoint.sh"]

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

EXPOSE 5555

CMD ["./server"]
