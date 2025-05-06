FROM golang:1.22.4 as builder
WORKDIR /go/src/app

ENV GO111MODULE=on

RUN go install github.com/cespare/reflex@latest

RUN go get github.com/google/uuid
RUN go get github.com/rabbitmq/amqp091-go
RUN go get github.com/ilyakaznacheev/cleanenv
RUN go get gopkg.in/telebot.v3

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download

COPY main-service/. main-service/.

RUN go build -o ./run ./main-service/cmd/.

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /go/src/app/run .

CMD ["./run"]
