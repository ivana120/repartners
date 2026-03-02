FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /pack-calculator ./cmd/server

FROM alpine:3.19

WORKDIR /app

RUN adduser -D -g '' appuser
COPY --from=builder /pack-calculator /app/pack-calculator
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENV PORT=8080
ENV PACK_SIZES="250,500,1000,2000,5000"

ENTRYPOINT ["/app/pack-calculator"]
