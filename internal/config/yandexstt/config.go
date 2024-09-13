package yandexstt

const (
	API_KEY        = "AQVNy3RDQQxrKiTsZJdB1ROjOMucYrAJk9f1KnT3" // API-Key для Yandex Speech Kit
	STORAGE_KEY_ID = "YCAJEfHe4j4TL9uMWwiAnCNNs"                // ID ключа для AWS
	STORAGE_KEY    = "YCP6_D3De-jmgeRmK8w7EzPTEjaTInV7GfvbDGlR" // Ключ AWS
	BUCKET_NAME    = "rvrecognition2"                           // Имя бакета в Yandex Object Storage

	SEND_RECOGNITION_URL  = "https://stt.api.cloud.yandex.net:443/stt/v3/recognizeFileAsync"
	CHECK_RECOGNITION_URL = "https://operation.api.cloud.yandex.net/operations/"
	GET_RECOGNITION_URL   = "https://stt.api.cloud.yandex.net:443/stt/v3/getRecognition?operation_id="
)
