FROM golang:1.16-bullseye
LABEL org.opencontainers.image.source = "https://github.com/skobaken7/docker-hosts";

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

CMD go run .
