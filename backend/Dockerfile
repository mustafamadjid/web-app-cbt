# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/web-app-cbt ./cmd/http

FROM golang:1.25-bookworm AS goose

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata wget \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p \
        /app/public/uploads/image/bank_soal \
        /app/public/uploads/document/import_soal \
        /app/public/uploads/document/pengumuman

COPY --from=builder /out/web-app-cbt /app/web-app-cbt
COPY --from=goose /go/bin/goose /usr/local/bin/goose
COPY internal/db/migrations /app/internal/db/migrations

ENV APP_DIR=/app
ENV UPLOAD_DIR=/app/public/uploads
ENV ENV=production

EXPOSE 8080

CMD ["/app/web-app-cbt"]
