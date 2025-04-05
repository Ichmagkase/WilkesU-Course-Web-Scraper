FROM golang:1.24-bookworm

COPY . .

RUN go mod download

RUN go build -o main

CMD ["./main"]
