package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"rvRecognitionService/structures"
	"strings"
	"time"
)

const API_KEY = "AQVNy3RDQQxrKiTsZJdB1ROjOMucYrAJk9f1KnT3"     // API-Key для Yandex Speech Kit
const STORAGE_KEY_ID = "YCAJEfHe4j4TL9uMWwiAnCNNs"             // ID ключа для AWS
const STORAGE_KEY = "YCP6_D3De-jmgeRmK8w7EzPTEjaTInV7GfvbDGlR" // Ключ AWS
const BUCKET_NAME = "rvrecognition2"                           // Имя бакета в Yandex Object Storage

func uploadToBucket(fullPath string) (string, error) {
	var bucketAddress string
	fparams := strings.Split(fullPath, "/")
	key := fparams[len(fparams)-1]

	log.Println("Uploading to bucket")

	cmd := exec.Command("aws", "--endpoint-url=https://storage.yandexcloud.net/", "s3", "cp", fullPath, "s3://"+BUCKET_NAME+"/"+key)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if err != nil {
		log.Println(fullPath, "-> ERROR WHILE UPLOADING FILE TO BUCKET:", err)
		log.Println(cmd)
		return "", err
	}

	log.Println(fullPath, "uploaded to bucket")

	bucketAddress = "https://storage.yandexcloud.net/" + BUCKET_NAME + "/" + key
	return bucketAddress, nil
}

func getAllSpeakerTexts(lines [][]byte) (string, error) {
	log.Println("Started parcer of responses from Yandex GPT")
	trueLines := make([][]byte, 0, len(lines))

	for _, line := range lines {
		if strings.Contains(string(line), "\"final\":{\"alternatives\":") {
			trueLines = append(trueLines, line)
		}
	}

	GetResponses := make([]structures.GetResponse, len(trueLines))
	for i, line := range trueLines {
		err := json.Unmarshal(line, &(GetResponses[i]))
		if err != nil {
			return "", err
		}
	}

	currentChannelTag := GetResponses[0].Result.ChannelTag
	speachParts := make([]string, 0, len(GetResponses))
	var currentChannelText string = ""

	for i, resp := range GetResponses {
		if resp.Result.ChannelTag != currentChannelTag {
			currentChannelTag = resp.Result.ChannelTag
			speachParts = append(speachParts, strings.TrimSpace(currentChannelText))
			currentChannelText = ""
		}
		currentChannelText += " "
		currentChannelText += resp.Result.Final.Alternatives[0].Text
		if i == (len(GetResponses) - 1) {
			speachParts = append(speachParts, strings.TrimSpace(currentChannelText))
		}
	}

	if len(speachParts) == 1 {
		return speachParts[0], nil
	}
	var text string = ""

	for _, part := range speachParts {
		text += "— " + strings.TrimSpace(part) + "\n\n"
	}

	text = strings.TrimSpace(text)

	return text, nil
}

func getResult(id string) (string, error) {
	log.Println("Started getResult for id:", id)
	checkHttpURL := "https://operation.api.cloud.yandex.net/operations/" + id
	getHttpURL := "https://stt.api.cloud.yandex.net:443/stt/v3/getRecognition?operation_id=" + id

	request, errorq := http.NewRequest("GET", checkHttpURL, nil)
	if errorq != nil {
		return "", errorq
	}

	request.Header.Set("Authorization", "Api-Key "+API_KEY)

	client := &http.Client{}
	var CheckResponse structures.CheckResponse

	// Проверка на готовность
	for {
		resp, err := client.Do(request)

		if err != nil {
			log.Println("Error while senging check request:", err)
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return "", errors.New("yandex: " + resp.Status)
		}
		checkData := make([]byte, 300)
		q, err := resp.Body.Read(checkData)
		if err != nil {
			log.Println("Error while reading check data:", err)
			return "", err
		}

		if resp.StatusCode != 200 {
			return "", errors.New("Yandex check server is unavailable.\nStatus: " + resp.Status + "\nResult: " + string(checkData[:q]))
		}

		err = json.Unmarshal(checkData[:q], &CheckResponse)
		if err != nil {
			log.Println("Error while unmarshalling check response:", err)
			log.Println("\n\n>>>>>", string(checkData))
			return "", err
		}
		if CheckResponse.Done {
			break
		} else {
			<-time.After(1 * time.Second)
			log.Println("Done:", CheckResponse.Done)
			continue
		}
	}

	request, errorq = http.NewRequest("GET", getHttpURL, nil)
	if errorq != nil {
		return "", errorq
	}
	request.Header.Set("Authorization", "Api-Key "+API_KEY)

	resp, err := client.Do(request)

	if err != nil {
		log.Println("Error while senging get response:", err)
		return "", err
	}

	reader := bufio.NewReader(resp.Body)

	var lines [][]byte

	var i = 0
	for {
		l, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		lines = append(lines, l)
		i++
	}

	result, err := getAllSpeakerTexts(lines)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return result, nil
}

func sendRequest(fileData string, lang string, isDialog bool) (string, error) {
	log.Println("Sending request to Yandex STT for", fileData)
	httpURL := "https://stt.api.cloud.yandex.net:443/stt/v3/recognizeFileAsync"
	httpMethod := "POST"

	splits := strings.Split(fileData, ".")
	typeF := splits[len(splits)-1]
	log.Println("File type:", typeF)
	// Инициализация тела запроса
	httpBody := &structures.SendRequest{}

	httpBody.URI = fileData                                                                          // Путь к файлу в бакете
	httpBody.RecognitionModel.Model = "general:rc"                                                   // Выбор модели: general / general:rc
	httpBody.RecognitionModel.AudioFormat.ContainerAudio.ContainerAudioType = strings.ToUpper(typeF) // Указание типа файла
	httpBody.RecognitionModel.LanguageRestriction.LanguageCode = []string{lang}                      // Указание языков
	httpBody.RecognitionModel.LanguageRestriction.RestrictionType = "WHITELIST"                      // Указание параметра использования языков
	httpBody.RecognitionModel.TextNormalization.TextNormalization = "TEXT_NORMALIZATION_DISABLED"    // Отключение нормализации текста
	if isDialog {
		httpBody.SpeakerLabeling.SpeakerLabeling = "SPEAKER_LABELING_ENABLED" // Подключение разделения на спикеров
	} else {
		httpBody.SpeakerLabeling.SpeakerLabeling = "SPEAKER_LABELING_DISABLED"
	}
	httpBody.RecognitionModel.AudioProcessingType = "FULL_DATA" // Выбор, как обрабатывать аудио: в реальном времени или после получения

	data, err := json.Marshal(httpBody)
	if err != nil {
		return "", err
	}

	request, errorq := http.NewRequest(httpMethod, httpURL, bytes.NewBuffer(data))

	if errorq != nil {
		return "", errorq
	}

	request.Header.Set("Authorization", "Api-Key "+API_KEY)
	request.Header.Set("x-data-logging-enabled", "true")

	client := &http.Client{}

	log.Println("Sending file to server")
	resp, err2 := client.Do(request)

	if err2 != nil {
		log.Println("Error while file to server")
		return "", err2
	}
	log.Println("File sent. Status:", resp.Status)
	defer resp.Body.Close()
	log.Println("Send file status:", resp.Status)
	var objBody structures.SendResponse

	byteBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error while reading body:", err)
		log.Println(string(byteBody))
		return "", err
	}

	err = json.Unmarshal(byteBody, &objBody)

	if err != nil {
		log.Println("Error while unmarshalling send response:", err)
		log.Println(string(byteBody))
		return "", err
	}

	log.Println(fileData, "->", resp)
	return objBody.Id, nil
}

// Функция распознавания речи в файле.
// Returns RawText, NormText, Error
func recognize(filePath string, lang string, dialog bool) (string, string, error) {
	log.Println("Startring recognition for", filePath)
	bucketFilePath, err := uploadToBucket(filePath)

	if err != nil {
		log.Println("Error while uploading file to bucket:", err)
		return "", "", err
	}

	id, err := sendRequest(bucketFilePath, lang, dialog)

	if err != nil {
		log.Println("Error while sending request:", err)
		return "", "", err
	}

	rawText, err := getResult(id) // Голая расшифровка

	if err != nil {
		log.Println("Error while getting result from Yandex:", err)
		return "", "", err
	}

	normalizedText := normilize(rawText) // Расшифровка, нормализованная через нейросеть

	return rawText, normalizedText, nil
}
