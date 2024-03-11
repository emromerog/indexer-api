# Golang image
FROM golang:1.22.1-alpine

# Env vars
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Project setup - copiamos todo el proyecto a /app
RUN mkdir app
COPY . /app
WORKDIR /app/cmd

# Build - el ejecutable se llamara main
RUN go build -o main .

# Port api
EXPOSE 8080

# Run - ejecutamos main
CMD [ "/app/cmd/main" ]