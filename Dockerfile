# Development Dockerfile with Air hot-reload
FROM golang:1.23-alpine

# Install Air for hot-reload
RUN go install github.com/air-verse/air@v1.61.7

WORKDIR /app

# Copy go.mod first for better caching
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

EXPOSE 80

# Run with Air for hot-reload
CMD ["air", "-c", ".air.toml"]
