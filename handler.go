package main

import (
	"bytes"
	"log"
	"net"
	"strings"
)

func handle(conn net.Conn) {

	// Чтение данных
	var data string = ""
	var lang string
	var filePath string

	for {
		byte_data := make([]byte, 50)
		_, err := conn.Read(byte_data)

		if err != nil {
			log.Println(conn.RemoteAddr(), "Error during reading content: ", err)
			conn.Close()
			return
		}

		trimmed_data := bytes.Trim(byte_data, "\x00")
		if string(trimmed_data[len(trimmed_data)-1]) == "\v" {
			data += string(trimmed_data[0 : len(trimmed_data)-1])
			break
		} else {
			data += string(trimmed_data)
		}
	}

	// Парминг данных
	data = strings.TrimSpace(data)
	lang = strings.TrimSpace(data[0:5])
	filePath = strings.TrimSpace(data[6:])

	log.Println(conn.RemoteAddr(), "Readed data: ", data)

	// Транскрибация
	log.Println(conn.RemoteAddr(), "Starting transcription")
	recognitionText := recognize(filePath, lang)

	// Возврат результата
	log.Println(conn.RemoteAddr(), "Sending data")
	_, err := conn.Write([]byte(recognitionText + "\v"))
	if err != nil {
		log.Println(conn.RemoteAddr(), "Error during sending data: ", err)
	}

	conn.Close()
}
