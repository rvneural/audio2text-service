package filedownloader

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var DOWNLOAD_SERVICE_URL = os.Getenv("FILE_DOWNLOADER")

type Response struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type Donwloader struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Donwloader {
	return &Donwloader{logger: logger}
}

func (d *Donwloader) Download(url string) (data []byte, fyleType string, fileName string, err error) {

	if url == "" {
		return nil, "", "", nil
	}

	full_url := DOWNLOAD_SERVICE_URL + url

	d.logger.Debug().Msgf("Downloading file from %s", full_url)
	resp, err := http.Get(full_url)
	if err != nil {
		d.logger.Error().Msgf("Error downloading file from %s: %s", full_url, err)
		return nil, "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.logger.Debug().Msgf("Error reading response body from %s: %s", full_url, err)
		return nil, "", "", err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		d.logger.Debug().Msgf("Error unmarshalling response from %s: %s", full_url, err)
		return nil, "", "", err
	}

	fileParts := strings.Split(response.Name, ".")
	fileType := fileParts[len(fileParts)-1]

	d.logger.Debug().Msgf("Downloaded file from %s: %s", full_url, response.Name)
	return response.Data, fileType, response.Name, nil
}
