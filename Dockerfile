FROM golang:1.16-bullseye

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

CMD go run .
