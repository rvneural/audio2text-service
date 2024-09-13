FROM golang:latest

LABEL authors="Artem Gafarov"
LABEL email="gafarov@realnoevremya.ru"
LABEL version="2.0"

WORKDIR /src/audio2text-service
ADD . .


ENV NORMALIZATION_ADDR="http://127.0.0.1:45680"
RUN go mod download && go mod verify

EXPOSE 80

WORKDIR /src/audio2text-service/cmd/audio2text-service
RUN go build

CMD ["./audio2text-service"]