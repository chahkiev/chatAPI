FROM golang:latest
LABEL maintainer="Chahkiev Magomed <chahkiev@mail.ru>"

WORKDIR /home/chatAPI

COPY go.mod go.sum ./

RUN go mod download

ADD . .

RUN go build -o main .

CMD ["./main"]