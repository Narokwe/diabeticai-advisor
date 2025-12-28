# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o diabeticai-advisor

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/diabeticai-advisor /app/diabeticai-advisor

EXPOSE 3400

CMD ["/app/diabeticai-advisor"]