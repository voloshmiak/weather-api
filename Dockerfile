FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && \
    go mod download

COPY . .

RUN echo "--- Contents of /app (after COPY . .) ---" && ls -R /app
RUN echo "--- Current directory (should be /app) ---" && pwd

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/weather-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]