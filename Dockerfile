# Dockerfile

# --- Build Stage ---
# Use an official Go image that satisfies the go.mod requirement (>= 1.24.3)
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies.
# This layer is cached by Docker, so dependencies are only re-downloaded
# if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application.
# -o tictactoe specifies the output file name.
# CGO_ENABLED=0 disables Cgo, which is needed for a static binary.
# -ldflags="-w -s" strips debugging information, reducing the binary size.
RUN CGO_ENABLED=0 go build -o /tictactoe -ldflags="-w -s" .

# --- Final Stage ---
# Use a minimal base image for the final container.
# 'scratch' is an empty image, providing the smallest possible size.
FROM scratch

# Set the working directory
WORKDIR /

# Copy the built binary from the 'builder' stage
COPY --from=builder /tictactoe /tictactoe

# Set the command to run when the container starts.
# This will execute our Tic-Tac-Toe game.
ENTRYPOINT ["/tictactoe"]
