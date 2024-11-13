package bucketuploader

import (
	"os"
	"strings"

	config "Audio2TextService/internal/config/yandexstt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/rs/zerolog"
)

const AWS_REGION = "ru-central1"

type AWSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	BucketName      string
	UploadTimeout   int
	BaseURL         string
}

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

	awsConfig := AWSConfig{
		AccessKeyID:     config.STORAGE_KEY_ID,
		AccessKeySecret: config.STORAGE_KEY,
		Region:          AWS_REGION,
		BucketName:      config.BUCKET_NAME,
		BaseURL:         "https://storage.yandexcloud.net/",
	}

	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AccessKeyID,
				awsConfig.AccessKeySecret,
				"",
			),
			Endpoint: aws.String(awsConfig.BaseURL),
		},
	))

	file, err := os.Open(filePath)
	if err != nil {
		u.Logger.Error().Msg("Error opening file: " + err.Error())
		return "", err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsConfig.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		u.Logger.Error().Msg("Error uploading file to Yandex Cloud Storage: " + err.Error())
		return "", err
	}
	// cmd := exec.Command("aws", "--endpoint-url=https://storage.yandexcloud.net/", "s3", "cp", filePath, "s3://"+config.BUCKET_NAME+"/"+key)
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// cmd.Stderr = &out
	// err := cmd.Run()

	// if err != nil {
	// 	u.Logger.Error().Msg("Error uploading file to Yandex Cloud Storage: " + out.String())
	// 	return "", err
	// }

	go os.Remove(filePath)

	bucketAddress = "https://storage.yandexcloud.net/" + config.BUCKET_NAME + "/" + key
	return bucketAddress, nil
}
