# Build stage
FROM golang:1.27-alpine3.24 AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o autodiscover .

# Runtime stage
FROM alpine:3.24

RUN apk --no-cache add ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder --chown=appuser:appgroup /app/autodiscover .

USER appuser

EXPOSE 8080

CMD ["./autodiscover"]
