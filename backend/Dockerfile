# ---- Build stage ----
FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod and go.sum first for better cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app (output binary to /app/server)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- Run stage ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/server .

# Copy .env if you want to bake it in (optional)
# COPY .env .env

EXPOSE 8080

CMD ["./server"]
