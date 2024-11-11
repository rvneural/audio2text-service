FROM debian:latest
LABEL maintainer="gafarov@realnoevremya.ru"
RUN apt-get update && apt-get upgrade
EXPOSE 8082
COPY . .
WORKDIR /build/linux
CMD [ "./audio2text-service" ]

