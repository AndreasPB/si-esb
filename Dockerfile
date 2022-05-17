FROM golang:alpine as base

WORKDIR /app

EXPOSE 9999

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o esb .

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
CMD air
