# Build Stage
FROM golang:1.21 AS build


# Copy sources inside the container
COPY . /app

WORKDIR /app

# Download dependencies
RUN go mod download && go mod verify

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o  main ./cmd/service

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main main

# Command to run the application
CMD ["./main"]