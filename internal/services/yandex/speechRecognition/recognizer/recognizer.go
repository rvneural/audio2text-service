package recognizer

import (
	stt "Audio2TextService/internal/models/json/yandexstt/send"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	config "Audio2TextService/internal/config/yandexstt"

	"github.com/rs/zerolog"
)

type Recognizer struct {
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Recognizer {
	return &Recognizer{Logger: logger}
}

func (r *Recognizer) SendRequest(path string, lang []string, dialog bool) (string, error) {

	r.Logger.Info().Msg("Sending request to Yandex Speech-to-Text for " + path)
	defer r.Logger.Info().Msg("Finished sending request to Yandex Speech-to-Text for " + path)

	httpURL := config.SEND_RECOGNITION_URL

	splits := strings.Split(path, ".")
	typeF := splits[len(splits)-1]

	// Инициализация тела запроса
	httpBody := r.getRequestBody(path, typeF, lang, dialog)

	data, err := json.Marshal(httpBody)
	if err != nil {
		r.Logger.Error().Msg("Error marshaling request body: " + err.Error())
		return "", err
	}

	resp, err2 := r.sendPostRequest(httpURL, data)

	if err2 != nil {
		return "", err2
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		r.Logger.Debug().Msgf("Error sending request body >>> status: %d", resp.StatusCode)
		return "", errors.New("yandex: " + resp.Status)
	}

	var objBody stt.Response

	byteBody, err := io.ReadAll(resp.Body)

	if err != nil {
		r.Logger.Error().Msg("Error reading response body: " + err.Error())
		return "", err
	}

	err = json.Unmarshal(byteBody, &objBody)

	if err != nil {
		r.Logger.Error().Msg("Error unmarshalling response body: " + err.Error())
		return "", err
	}

	return objBody.Id, nil
}

func (r *Recognizer) GetResponse(id string) ([][]byte, error) {

	r.Logger.Info().Msg("Getting response from Yandex Speech-to-Text for " + id)
	defer r.Logger.Info().Msg("Finished getting response from Yandex Speech-to-Text for " + id)

	getHttpURL := config.GET_RECOGNITION_URL + id

	done, err := r.waitForRecognition(id)

	if err != nil {
		r.Logger.Error().Msg("Error waiting for recognition: " + err.Error())
		return nil, err
	}
	if !done {
		r.Logger.Debug().Msg("Recognition is not done yet for ID: " + id)
		return nil, errors.New("recognition is not done")
	}

	request, errorq := http.NewRequest("GET", getHttpURL, nil)
	if errorq != nil {
		r.Logger.Error().Msg("Error creating GET request: " + errorq.Error())
		return nil, errorq
	}
	request.Header.Set("Authorization", "Api-Key "+config.API_KEY)

	client := &http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		r.Logger.Error().Msg("Error sending GET request: " + err.Error())
		return nil, err
	}

	reader := bufio.NewReader(resp.Body)

	var lines [][]byte

	var i = 0
	for {
		l, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			r.Logger.Error().Msg("Error reading response body: " + err.Error())
			return nil, err
		}
		lines = append(lines, l)
		i++
	}

	return lines, nil
}
