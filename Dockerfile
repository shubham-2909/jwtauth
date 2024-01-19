# Use a smaller base image for the final stage
FROM golang:alpine as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your application runs on
EXPOSE 8000

# Create a non-root user for running the application
RUN adduser -D -u 1001 appuser
USER appuser

# Command to run the executable
CMD ["./main"]