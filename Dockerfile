# Build Stage
FROM golang:1.21


# Copy sources inside the container
COPY . /app

WORKDIR /app

# Download dependencies
RUN go mod download && go mod verify

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o ./main ./cmd/service

# todo: application don't run with enabled cgo in slim alpine image, investigate how to run binary in slim image without golang