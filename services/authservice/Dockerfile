FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates 

WORKDIR /app/authservice

COPY common /common

COPY services/authservice/ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
        -o /auth \
        ./cmd

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /auth /auth 

ENTRYPOINT ["/auth"]