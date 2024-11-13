package yandexstt

import "os"

const (
	SEND_RECOGNITION_URL  = "https://stt.api.cloud.yandex.net:443/stt/v3/recognizeFileAsync"
	CHECK_RECOGNITION_URL = "https://operation.api.cloud.yandex.net/operations/"
	GET_RECOGNITION_URL   = "https://stt.api.cloud.yandex.net:443/stt/v3/getRecognition?operation_id="
)

var (
	API_KEY        = os.Getenv("API_KEY")        // API-Key для Yandex Speech Kit
	STORAGE_KEY_ID = os.Getenv("STORAGE_KEY_ID") // ID ключа для AWS
	STORAGE_KEY    = os.Getenv("STORAGE_KEY")    // Ключ AWS
	BUCKET_NAME    = os.Getenv("BUCKET_NAME")    // Имя бакета в Yandex Object Storage
)
