# Use the official Go image as a base
FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Install templ binary executables
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate the .go files from .templ files
RUN templ generate

# Build the Go app
RUN go build -o mulhall cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./mulhall"]
