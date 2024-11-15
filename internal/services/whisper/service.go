package whisper

import (
	"Audio2TextService/internal/models/json/whisper"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

type WhisperTransctiptor struct {
	URL    string
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *WhisperTransctiptor {
	return &WhisperTransctiptor{
		Logger: logger,
		URL:    os.Getenv("WHISPER_URL"),
	}
}

func (w *WhisperTransctiptor) RecognizeAudio(fileData []byte, fileType string) (string, error) {
	if w.URL == "" {
		return "", fmt.Errorf("WHISPER_URL is not set")
	}

	request := whisper.Request{
		FileData: fileData,
		FileName: w.getRandonName(30) + "." + fileType,
	}

	byteRequest, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	response, err := http.Post(w.URL, "application/json", bytes.NewBuffer(byteRequest))
	if err != nil {
		return "", err
	}

	byteData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseData whisper.Response
	err = json.Unmarshal(byteData, &responseData)
	return responseData.Text, err
}

func (w *WhisperTransctiptor) getRandonName(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	rune_name := make([]rune, length)
	for i := range rune_name {
		rune_name[i] = letters[rand.Intn(len(letters))]
	}
	return string(rune_name)
}
