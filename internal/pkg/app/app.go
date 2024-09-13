package app

import (
	endpoint "Audio2TextService/internal/endpoint/app"
	service "Audio2TextService/internal/services"

	handler "Audio2TextService/internal/transport/rest"

	"Audio2TextService/internal/services/yandex/normalization"
	speechRecognizer "Audio2TextService/internal/services/yandex/speechRecognition"

	uploader "Audio2TextService/internal/services/yandex/speechRecognition/bucketUploader"
	"Audio2TextService/internal/services/yandex/speechRecognition/parser"
	"Audio2TextService/internal/services/yandex/speechRecognition/recognizer"

	processor "Audio2TextService/pkg/fileprocessor"
	"Audio2TextService/pkg/fileprocessor/converter"

	"github.com/rs/zerolog"
)

type App struct {
	endpoint *endpoint.Endpoint
	service  *service.Service
	handler  *handler.Audio2TextHandler
	logger   *zerolog.Logger
}

func New(logger *zerolog.Logger) *App {

	// Инициализация сервисов и обработчика
	uploader := uploader.New(logger)
	recognizer := recognizer.New(logger)
	parser := parser.New(logger)

	speechRecognition := speechRecognizer.New(uploader, recognizer, parser, logger)
	normalization := normalization.New(logger)

	converter := converter.New(logger)
	processor := processor.New(converter, logger)

	service := service.New(speechRecognition, normalization, processor, logger)

	handler := handler.New(service, logger)
	endpoint := endpoint.New(handler, logger)

	return &App{endpoint: endpoint, handler: handler, service: service, logger: logger}
}

func (a *App) Run() error {
	// Запуск сервера
	a.logger.Info().Msg("Starting Audio2TextService...")
	defer a.logger.Info().Msg("Audio2TextService stopped")
	return a.endpoint.Start()
}
