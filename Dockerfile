# Use an official Golang image as the base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project to the container
COPY . .

# Build the Go app - output the binary as /app/main
RUN go build -o main ./server/main.go

# Expose the application port (if needed)
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
