package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// Адрес для обращения к сервису нормализации текста
const addr string = "http://127.0.0.1:45678/text/normalize"

// Структура для отправки запроса к сервису нормализации
type Request struct {
	Text string `json:"text"` // Текст для нормализации
}

// Структура для получения ответа от сервиса нормализации
type Response struct {
	OldText string `json:"oldText"` // Исходный текст
	NewText string `json:"newText"` // Нормализованный текст
}

// Основная функция для нормализации текста
// Принимает текст на вход и возвращает нормализованный текст
func normalize(text string) string {
	var normalizedText string = ""

	// Создаем структуру для отправки запроса
	request := Request{Text: text}

	// Маршаллируем структуру в JSON-формате для отправки
	bytesRequest, err := json.Marshal(request)

	if err != nil {
		log.Println("Error marshalling request", err)
		return err.Error()
	}

	// Отправляем запрос к сервису нормализации
	resp, err := http.Post(addr, "application/json", bytes.NewReader(bytesRequest))
	if err != nil {
		log.Println("Error while sending request:", err)
		return err.Error()
	}
	defer resp.Body.Close()

	// Получаем ответ от сервиса
	var response Response
	byteResponse, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading response body:", err)
		return err.Error()
	}

	// Выводим полученный ответ на консоль
	log.Println(strings.TrimSpace(string(byteResponse)))

	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return err.Error()
	}

	normalizedText = response.NewText

	log.Printf("Normalized text: %s\n", strings.ReplaceAll(normalizedText, "\n", " "))

	return normalizedText
}
