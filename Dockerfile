# Use an official, lightweight Go image as a base
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first to leverage Docker's build cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application's source code
COPY . .

# Build the Go application into a single executable named "server"
# CGO_ENABLED=0 is a standard practice for creating a static binary
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/bot

# This is the command that will run when the container starts
CMD ["/app/server"]