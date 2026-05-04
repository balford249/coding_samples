FROM golang:1.22 AS builder

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    gcc \
    g++


WORKDIR /src
COPY . .



# Build API
RUN go build -o eval-api ./api

# Build Consumer
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#    go build -o consumer ./cmd/consumer

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /src/eval-api .
COPY --from=builder /src/api/config.json .
#COPY --from=builder /src/consumer .
