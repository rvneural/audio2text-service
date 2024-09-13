package speechRecognition

import (
	"fmt"

	"github.com/rs/zerolog"
)

type SpeechRecognition struct {
	Uploader   Uploader
	Recognizer Recognizer
	Parser     Parser
	Logger     *zerolog.Logger
}

func New(uploader Uploader, recognizer Recognizer, parser Parser, logger *zerolog.Logger) *SpeechRecognition {
	return &SpeechRecognition{Uploader: uploader, Recognizer: recognizer, Parser: parser, Logger: logger}
}

func (s *SpeechRecognition) RecognizeAudio(filePath string, lang []string, dialog bool, uniqPhraseSplitter string, maxLength int) ([]string, error) {

	logLine := fmt.Sprintf("New Speech Recognizer request for %s", filePath)
	s.Logger.Info().Msg(logLine)
	defer s.Logger.Info().Msg(fmt.Sprintf("Finished Speech Recognizer request for %s", filePath))

	bucketFilePath, err := s.Uploader.Upload(filePath) // Загрузка файла в бакет

	if err != nil {
		s.Logger.Error().Msg("Error uploading file to bucket: " + err.Error())
		return nil, err
	}

	id, err := s.Recognizer.SendRequest(bucketFilePath, lang, dialog) // Отправка запроса на распознавание

	if err != nil {
		s.Logger.Error().Msg("Error sending request to Speech-to-Text: " + err.Error())
		return nil, err
	}

	lines, err := s.Recognizer.GetResponse(id)

	if err != nil {
		s.Logger.Error().Msg("Error getting response from Speech-to-Text: " + err.Error())
		return nil, err
	}

	rawText, err := s.Parser.Parse(lines, uniqPhraseSplitter, maxLength)

	if err != nil {
		s.Logger.Error().Msg("Error parsing response: " + err.Error())
		return nil, err
	}

	return rawText, nil
}
