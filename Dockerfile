FROM golang:1.24-bookworm

COPY go.mod go.sum /build/

WORKDIR /build

RUN go mod download

COPY ./src/ /build/

RUN go build -o main

CMD ["./main"]

