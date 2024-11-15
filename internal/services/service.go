package services

import (
	config "Audio2TextService/internal/config/service"
	"strings"

	"github.com/rs/zerolog"
)

type Service struct {
	Recognition        Recognition
	WhisperRecognition WhisperRecognition
	Normalization      Normalization
	FileProcessor      FileProcessor
	Logger             *zerolog.Logger
}

func New(recognition Recognition, normalization Normalization, processor FileProcessor, logger *zerolog.Logger) *Service {
	return &Service{Recognition: recognition, Normalization: normalization, Logger: logger, FileProcessor: processor}
}

// Конвертирует аудиофайл в текст с указанием
func (s *Service) ConvertAudioToText(fileData []byte, fileType string, model string, lang []string, dialog bool) (rawText string, normText string, err error) {
	fileType = strings.ToLower(fileType)
	s.Logger.Info().Msg("Service: Converting audio to text")
	s.Logger.Info().Msg("Service: Starting file processing")
	// Обрабатываем файл и получаем путь до него

	if model == "whisper" {
		s.Logger.Info().Msg("Using Whisper model")
		transcriptedText, err := s.WhisperRecognition.RecognizeAudio(fileData, fileType)
		return transcriptedText, transcriptedText, err
	}

	filePath, err := s.FileProcessor.ProcessFile(fileData, fileType)

	if err != nil {
		s.Logger.Error().Msg("Error processing file: " + err.Error())
		return "", "", err
	}
	s.Logger.Debug().Msg("File processed OK")

	s.Logger.Info().Msg("Converting " + filePath + " to text")
	defer s.Logger.Info().Msg("Finished converting " + filePath + " to text")

	// Распознаем аудиофайл и получаем сырой текст
	rawLines, err := s.Recognition.RecognizeAudio(filePath, lang, dialog, config.UNIQ_PHRASE_SPLITTER, config.MAX_LENGTH)

	if err != nil {
		s.Logger.Error().Msg("Error recognizing audio: " + err.Error())
		return "", "", err
	}

	// Нормализуем текст - расставляем знаки препинания и заглавные буквы каждого предложения
	if len(rawLines) == 1 {
		return rawLines[0], s.Normalization.NormalizeText(rawLines[0]), nil
	} else {
		var dialogStart = "—"
		var lineSplitter = "\n\n"

		var rawText string = ""
		var normText string = ""

		for _, line := range rawLines {
			if !dialog {
				line = strings.TrimSpace(strings.TrimPrefix(line, config.UNIQ_PHRASE_SPLITTER))
			}
			var isStartPhrase = strings.HasPrefix(line, config.UNIQ_PHRASE_SPLITTER)
			p_line := strings.TrimSpace(strings.TrimPrefix(line, config.UNIQ_PHRASE_SPLITTER))
			rawText += line + lineSplitter
			normLine := s.Normalization.NormalizeText(p_line)
			if isStartPhrase {
				normText += dialogStart + " "
			}
			normText += normLine + lineSplitter
		}
		return rawText, strings.TrimSpace(normText), nil
	}
}
