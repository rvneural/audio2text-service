package main

import (
	"log"
	"net"
	"time"
)

func startServer() {
	var addr string = "127.0.0.1:45679"
	// Начинаем прослушивать сервер
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic("Ошибка чтения адреса: ", err)
	}

	log.Println("Waiting for requests: " + addr)

	defer server.Close()

	for {
		// Проверяем соединения каждые 300 миллисекунд
		<-time.After(300 * time.Millisecond)
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Вызываем обработчик запроса
		log.Println(conn.RemoteAddr(), "Started handle")
		go handle(conn)
	}
}
