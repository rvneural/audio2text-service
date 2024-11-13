package recognizer

import (
	model "Audio2TextService/internal/models/json/yandexstt/check"
	stt "Audio2TextService/internal/models/json/yandexstt/send"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	config "Audio2TextService/internal/config/yandexstt"
)

func (r *Recognizer) waitForRecognition(id string) (bool, error) {
	checkHttpURL := config.CHECK_RECOGNITION_URL + id

	request, errorq := http.NewRequest("GET", checkHttpURL, nil)
	if errorq != nil {
		return false, errorq
	}

	request.Header.Set("Authorization", "Api-Key "+config.API_KEY)

	client := &http.Client{}
	var CheckResponse model.Response

	for {
		resp, err := client.Do(request)

		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return false, errors.New("yandex: " + resp.Status)
		}
		checkData := make([]byte, 300)
		q, err := resp.Body.Read(checkData)
		if err != nil {
			return false, err
		}

		if resp.StatusCode != 200 {
			return false, errors.New("Yandex check server is unavailable.\nStatus: " + resp.Status + "\nResult: " + string(checkData[:q]))
		}

		err = json.Unmarshal(checkData[:q], &CheckResponse)
		if err != nil {
			return false, err
		}
		if strings.ToLower(os.Getenv("DEBUG_MODE")) == "true" {
			r.Logger.Debug().Msgf(string(checkData))
		}
		if CheckResponse.Done {
			return true, nil
		} else {
			<-time.After(1 * time.Second)
			continue
		}
	}
}

func (r *Recognizer) getRequestBody(fileURI, typeF string, lang []string, isDialog bool) stt.Request {

	var httpBody stt.Request

	httpBody.URI = fileURI                                                                           // Путь к файлу в бакете
	httpBody.RecognitionModel.Model = "general:rc"                                                   // Выбор модели: general / general:rc
	httpBody.RecognitionModel.AudioFormat.ContainerAudio.ContainerAudioType = strings.ToUpper(typeF) // Указание типа файла
	httpBody.RecognitionModel.LanguageRestriction.LanguageCode = lang                                // Указание языков
	httpBody.RecognitionModel.LanguageRestriction.RestrictionType = "WHITELIST"                      // Указание параметра использования языков
	httpBody.RecognitionModel.TextNormalization.TextNormalization = "TEXT_NORMALIZATION_DISABLED"    // Отключение нормализации текста
	if isDialog {
		httpBody.SpeakerLabeling.SpeakerLabeling = "SPEAKER_LABELING_ENABLED" // Подключение разделения на спикеров
	} else {
		httpBody.SpeakerLabeling.SpeakerLabeling = "SPEAKER_LABELING_DISABLED"
	}
	httpBody.RecognitionModel.AudioProcessingType = "FULL_DATA" // Выбор, как обрабатывать аудио: в реальном времени или после получения
	if os.Getenv("DEBUG_MODE") == "true" {
		r.Logger.Debug().Msgf("Request body: %v", httpBody)
	}
	return httpBody
}

func (r *Recognizer) sendPostRequest(url string, data []byte) (*http.Response, error) {
	request, errorq := http.NewRequest("POST", url, bytes.NewBuffer(data))

	if errorq != nil {
		return nil, errorq
	}

	request.Header.Set("Authorization", "Api-Key "+config.API_KEY)
	request.Header.Set("x-data-logging-enabled", "true")

	client := &http.Client{}

	return client.Do(request)
}
