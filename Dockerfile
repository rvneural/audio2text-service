FROM debian:latest
LABEL maintainer="gafarov@realnoevremya.ru"
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y ca-certificates
RUN apt install ffmpeg -y
EXPOSE 8082
COPY . .
WORKDIR /build/linux
CMD [ "./audio2text-service" ]

