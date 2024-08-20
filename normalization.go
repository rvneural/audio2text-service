package main

import (
	"bytes"
	"log"
	"net"
	"strings"
)

const addr string = "127.0.0.1:45680"

func normilize(text string) string {
	var normalizedText string = ""
	var request = text + "\v"

	log.Println("Connecting to normalization service")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Normalization service is not available", err)
		return text
	}

	defer conn.Close()

	_, err = conn.Write([]byte(request + "\v"))
	if err != nil {
		log.Println("Normalization service is not available", err)
		return text
	}

	batch_size := len(text) + 50
	for {
		text_sile := make([]byte, batch_size)
		_, err = conn.Read(text_sile)
		if err != nil {
			log.Println("Error during reading transcribed text: ", err)
			return text
		}
		f_slice := bytes.Trim(text_sile, "\x00")
		if string(f_slice[len(f_slice)-1]) == "\v" {
			normalizedText += string(f_slice[:len(f_slice)-1])
			break
		}
		normalizedText += string(f_slice)
	}

	return strings.TrimSpace(normalizedText)
}
