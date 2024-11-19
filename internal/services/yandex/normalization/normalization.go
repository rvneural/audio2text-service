package normalization

import (
	server "Audio2TextService/internal/config/normalization"
	"Audio2TextService/internal/models/json/normalization"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type Normalization struct {
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Normalization {
	return &Normalization{Logger: logger}
}

// Нормализует текст - расставляет знаки препинания и заглавные буквы
func (n *Normalization) NormalizeText(text string) string {

	n.Logger.Info().Msg("Normalizing text")
	defer n.Logger.Info().Msg("Finished normalizing text")

	var normalizedText string = ""

	// Создаем структуру для отправки запроса
	request := normalization.Request{
		Text:   text,
		Model:  "pro",
		Prompt: "{{ normalize }}",
	}

	// Маршаллируем структуру в JSON-формате для отправки
	bytesRequest, err := json.Marshal(request)

	if err != nil {
		n.Logger.Error().Msgf("Error marshalling request %s", err.Error())
		return err.Error()
	}

	httpRequest, err := http.NewRequest("POST", server.ADDR, bytes.NewReader(bytesRequest))
	if err != nil {
		n.Logger.Error().Msgf("Error creating POST request %s", err.Error())
		return err.Error()
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(httpRequest)
	if err != nil {
		n.Logger.Error().Msgf("Error sending POST request %s", err.Error())
		return err.Error()
	}
	defer resp.Body.Close()

	// Получаем ответ от сервиса
	var response normalization.Response
	byteResponse, err := io.ReadAll(resp.Body)

	if err != nil {
		n.Logger.Error().Msgf("Error reading response data %s", err.Error())
		return err.Error()
	}

	err = json.Unmarshal(byteResponse, &response)
	if err != nil {
		n.Logger.Error().Msgf("Error unmarshalling response %s", err.Error())
		return err.Error()
	}

	normalizedText = response.NewText

	return normalizedText
}
