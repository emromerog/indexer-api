# Golang image
FROM golang:1.22.1-alpine

# Env vars
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

## Project setup

# Create app directory inside container
RUN mkdir app

# cd /app
WORKDIR /app

COPY go.mod .

# Download all dependencies
RUN go mod download

# Copies all files from local to container
COPY . /app
#COPY . .

# Build Golang API
RUN go build -o ./cmd/main ./cmd

# Expose port 8080 the indexer API in golang
EXPOSE 8080

WORKDIR /cmd

# Run
CMD ["./main"]
# CMD ["./cmd/main"]
