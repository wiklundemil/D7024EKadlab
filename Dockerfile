# Stage 1: Build the Go binary
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies first.
COPY go.mod ./
RUN go mod download

# Now copy the rest of the application code.
COPY . .

# Build the application. 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kademlia_app main.go

# Stage 2: Create a small runtime image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/kademlia_app .

# Make sure the binary is executable. add root 
RUN chmod +x /root/kademlia_app

# Expose the port (3000, for example)
EXPOSE 3000

# Run the application
CMD ["./kademlia_app"]
