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

	var need_normalization bool = true

	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]

	need_normalization = strings.Contains(fileName, "norm_")
	log.Println("Normalization:", need_normalization)

	log.Println(conn.RemoteAddr(), "Readed data: ", data)
	log.Println("Language:", lang)
	log.Println("File path:", filePath)

	// Транскрибация
	log.Println(conn.RemoteAddr(), "Starting transcription")
	recognitionText, err_rec := recognize(filePath, lang)

	if err_rec == nil {
		log.Println("Recognition text:", recognitionText)
		var recognitionNormText string = recognitionText
		if recognitionText != "" && need_normalization {
			// Нормализация результата
			recognitionNormText = normilize(recognitionText)
		}

		if recognitionNormText == "" {
			recognitionNormText = "Yandex Server in unavailable"
		}
		log.Println("Recognition text:", recognitionText)

		// Возврат результата
		log.Println(conn.RemoteAddr(), "Sending data to main server: ", recognitionNormText)
		_, err := conn.Write([]byte(recognitionNormText + "\v"))
		if err != nil {
			log.Println(conn.RemoteAddr(), "Error during sending data: ", err)
		}
	} else {
		_, err := conn.Write([]byte(err_rec.Error() + "\v"))
		if err != nil {
			log.Println(conn.RemoteAddr(), "Error during sending data: ", err)
		}
	}

	conn.Close()
}
