FROM golang:1.21-alpine

WORKDIR /app

COPY . /app

RUN go build -o test

CMD ["./test"]