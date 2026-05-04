FROM golang:1.22 AS builder

WORKDIR /src
COPY . .

# Build API
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o eval-api ./api

# Build Consumer
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#    go build -o consumer ./cmd/consumer

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /src/api .
#COPY --from=builder /src/consumer .

CMD ["./api"]