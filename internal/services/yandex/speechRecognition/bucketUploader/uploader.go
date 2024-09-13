package bucketuploader

import (
	"bytes"
	"os/exec"
	"strings"

	config "Audio2TextService/internal/config/yandexstt"

	"github.com/rs/zerolog"
)

type Uploader struct {
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Uploader {
	return &Uploader{Logger: logger}
}

// Загружает аудиофайл в Yandex Cloud Storage и возвращает адрес загруженного файла
func (u *Uploader) Upload(filePath string) (string, error) {

	u.Logger.Info().Msg("Uploading file to Yandex Cloud Storage: " + filePath)
	defer u.Logger.Info().Msg("Finished uploading file to Yandex Cloud Storage: " + filePath)

	var bucketAddress string
	fparams := strings.Split(filePath, "/")
	key := fparams[len(fparams)-1]

	cmd := exec.Command("aws", "--endpoint-url=https://storage.yandexcloud.net/", "s3", "cp", filePath, "s3://"+config.BUCKET_NAME+"/"+key)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if err != nil {
		u.Logger.Error().Msg("Error uploading file to Yandex Cloud Storage: " + out.String())
		return "", err
	}

	bucketAddress = "https://storage.yandexcloud.net/" + config.BUCKET_NAME + "/" + key
	return bucketAddress, nil
}
