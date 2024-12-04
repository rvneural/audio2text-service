package rest

import (
	"Audio2TextService/internal/models/json/client"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Service interface {
	ConvertAudioToText(fileData []byte, fileType string, model string, lang []string, dialog bool) (rawText string, normText string, err error)
}

type Donwloader interface {
	Download(url string) ([]byte, string, string, error)
}

type DBWorker interface {
	RegisterOperation(uniqID string, operation_type string, user_id int) error
	SetResult(uniqID string, data []byte) error
}

type Audio2TextHandler struct {
	service    Service
	downloader Donwloader
	logger     *zerolog.Logger
	dbworker   DBWorker
}

func New(service Service, downloader Donwloader, dbworker DBWorker, logger *zerolog.Logger) *Audio2TextHandler {
	return &Audio2TextHandler{service: service, downloader: downloader, logger: logger, dbworker: dbworker}
}

func (h *Audio2TextHandler) HandleRequest(c echo.Context) error {

	h.logger.Info().Msgf("Received request: %+v", c.RealIP())
	defer h.logger.Info().Msgf("Received response: %+v", c.RealIP())

	var request = new(client.Request)

	if c.Request().Header.Get("Content-Type") == "" {
		h.logger.Error().Msg("Missing content type header")
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Missing content type header",
			Details: "Content type header is required\nUser application/json or application/x-www-form-urlencoded or application/xml"})
	}

	err := c.Bind(request)

	if os.Getenv("DEBUG_MODE") == "true" {
		h.logger.Info().Msgf("Request dialog: %+v", request.Dialog)
		h.logger.Info().Msgf("Request lang: %+v", request.Languages)
	}

	if err != nil {
		h.logger.Error().Msg("Error binding request body: " + err.Error())
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Invalid request body",
			Details: err.Error()})
	}

	if request.URL != "" && len(request.File.Data) != 0 {
		h.logger.Error().Msg("Both file data and url are specified")
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Both file data and url are specified",
			Details: "Only one of them should be specified"})
	}

	if request.URL != "" {
		h.logger.Info().Msg("Downloading file from URL: " + request.URL)
		request.File.Data, request.File.Type, request.File.Name, err = h.downloader.Download(request.URL)

		if err != nil {
			h.logger.Error().Msg("Error downloading file: " + err.Error())
			return c.JSON(http.StatusInternalServerError, client.Error{Error: "Error downloading file",
				Details: err.Error()})
		}
	}

	if request.Operation_ID != "" {
		id, err := strconv.Atoi(request.UserID)
		if err != nil {
			id = 0
		}
		go h.dbworker.RegisterOperation(request.Operation_ID, "audio", id)
	}

	if len(request.File.Data) == 0 || request.File.Type == "" {
		h.logger.Error().Msg("Missing file data or file type")
		return c.JSON(http.StatusBadRequest, client.Error{Error: "Missing file data or file type",
			Details: "File data and file type are required"})
	}

	h.logger.Info().Msg("Converting file")
	rawText, normText, err := h.service.ConvertAudioToText(request.File.Data,
		request.File.Type, request.Model, request.Languages, request.Dialog)

	h.logger.Info().Msg("File converted")

	if err != nil {
		h.logger.Error().Msg("Error converting audio to text: " + err.Error())
		return c.JSON(http.StatusInternalServerError, client.Error{Error: "Error converting audio to text",
			Details: err.Error()})
	}

	var response = client.Response{
		RawText:  rawText,
		NormText: normText,
	}

	if request.Operation_ID != "" {
		if request.File.Name == "" {
			request.File.Name = "recognized-file." + request.File.Type
		}
		dbResult := client.DBResult{
			FileName: request.File.Name,
			RawText:  rawText,
			NormText: normText,
		}
		data_result, _ := json.Marshal(dbResult)
		go h.dbworker.SetResult(request.Operation_ID, data_result)
	}
	return c.JSON(http.StatusOK, response)
}
