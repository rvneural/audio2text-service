package rest

import (
	client2 "Audio2TextService/internal/models/json/client"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Service interface {
	ConvertAudioToText(fileData []byte, fileType string, lang []string, dialog bool) (rawText string, normText string, err error)
}

type Donwloader interface {
	Download(url string) ([]byte, string, error)
}

type Audio2TextHandler struct {
	service    Service
	downloader Donwloader
	logger     *zerolog.Logger
}

func New(service Service, downloader Donwloader, logger *zerolog.Logger) *Audio2TextHandler {
	return &Audio2TextHandler{service: service, downloader: downloader, logger: logger}
}

func (h *Audio2TextHandler) HandleRequest(c echo.Context) error {

	h.logger.Info().Msgf("Received request: %+v", c.RealIP())
	defer h.logger.Info().Msgf("Received response: %+v", c.RealIP())

	var request = new(client2.Request)

	if c.Request().Header.Get("Content-Type") == "" {
		h.logger.Error().Msg("Missing content type header")
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Missing content type header",
			Details: "Content type header is required\nUser application/json or application/x-www-form-urlencoded or application/xml"})
	}

	h.logger.Info().Msgf("Binding request: %+v", c.Request)
	err := c.Bind(request)

	if err != nil {
		h.logger.Error().Msg("Error binding request body: " + err.Error())
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Invalid request body",
			Details: err.Error()})
	}

	if request.URL != "" && len(request.FileData) != 0 {
		h.logger.Error().Msg("Both file data and url are specified")
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Both file data and url are specified",
			Details: "Only one of them should be specified"})
	}

	if request.URL != "" {
		h.logger.Info().Msg("Downloading file from URL: " + request.URL)
		request.FileData, request.FileType, err = h.downloader.Download(request.URL)

		if err != nil {
			h.logger.Error().Msg("Error downloading file: " + err.Error())
			return c.JSON(http.StatusInternalServerError, client2.Error{Error: "Error downloading file",
				Details: err.Error()})
		}
	}

	if len(request.FileData) == 0 || request.FileType == "" {
		h.logger.Error().Msg("Missing file data or file type")
		return c.JSON(http.StatusBadRequest, client2.Error{Error: "Missing file data or file type",
			Details: "File data and file type are required"})
	}

	h.logger.Info().Msg("Converting file")
	rawText, normText, err := h.service.ConvertAudioToText(request.FileData,
		request.FileType, request.Languages, request.Dialog)

	h.logger.Info().Msg("File converted")
	h.logger.Info().Msg("Raw text: " + rawText)
	h.logger.Info().Msg("Normalized text: " + normText)

	if err != nil {
		h.logger.Error().Msg("Error converting audio to text: " + err.Error())
		return c.JSON(http.StatusInternalServerError, client2.Error{Error: "Error converting audio to text",
			Details: err.Error()})
	}

	var response = client2.Response{
		RawText:  rawText,
		NormText: normText,
	}
	return c.JSON(http.StatusOK, response)
}
