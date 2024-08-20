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

const API_KEY = "AQVNy3RDQQxrKiTsZJdB1ROjOMucYrAJk9f1KnT3"
const STORAGE_KEY_ID = "YCAJEfHe4j4TL9uMWwiAnCNNs"
const STORAGE_KEY = "YCP6_D3De-jmgeRmK8w7EzPTEjaTInV7GfvbDGlR"
const BUCKET_NAME = "rvrecognition2"

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
		log.Println("ERROR WHILE UPLOADING FILE TO BUCKET:", err)
		log.Println(cmd)
		return "", err
	}

	bucketAddress = "https://storage.yandexcloud.net/" + BUCKET_NAME + "/" + key
	return bucketAddress, nil
}

func getResult(id string) (string, error) {
	checkHttpURL := "https://operation.api.cloud.yandex.net/operations/" + id
	getHttpURL := "https://stt.api.cloud.yandex.net:443/stt/v3/getRecognition?operation_id=" + id

	request, errorq := http.NewRequest("GET", checkHttpURL, nil)
	if errorq != nil {
		return "", errorq
	}

	request.Header.Set("Authorization", "Api-Key "+API_KEY)

	client := &http.Client{}
	var GetResponse structures.GetResponse
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

	line, err := reader.ReadBytes('\n')

	var i = 0
	for i < 2 {
		l, _ := reader.ReadBytes('\n')
		log.Println(string(l))
		i++
	}

	reader.Reset(resp.Body)
	size := reader.Size()
	data := make([]byte, size)
	reader.Read(data)
	log.Println("\n\n\n\nDATA\n\n\n\n", string(data))

	if err != nil {
		log.Println("Error while reading get response:", err)
		log.Println(resp)
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("Yandex result server is unavailable.\nStatus: " + resp.Status + "\nResult: " + string(line))
	}

	err = json.Unmarshal(line, &GetResponse)

	if err != nil {
		log.Println("\n\n\n\n\nError while unmarshalling get response:", err)
		log.Println(string(line))
		return "", err
	}

	log.Println("Body:", GetResponse)
	return GetResponse.Result.Final.Alternatives[0].Text, nil
}

func sendRequest(fileData string, lang string) (string, error) {
	httpURL := "https://stt.api.cloud.yandex.net:443/stt/v3/recognizeFileAsync"
	httpMethod := "POST"

	splits := strings.Split(fileData, ".")
	typeF := splits[len(splits)-1]
	log.Println("File type:", typeF)
	// Инициализация тела запроса
	httpBody := &structures.SendRequest{}

	httpBody.URI = fileData
	httpBody.RecognitionModel.LanguageRestriction.LanguageCode = []string{lang}
	httpBody.RecognitionModel.Model = "general:rc"
	httpBody.RecognitionModel.AudioFormat.ContainerAudio.ContainerAudioType = strings.ToUpper(typeF)
	httpBody.RecognitionModel.LanguageRestriction.RestrictionType = "WHITELIST"

	data, err := json.Marshal(httpBody)
	if err != nil {
		return "", err
	}

	request, errorq := http.NewRequest(httpMethod, httpURL, bytes.NewBuffer(data))

	if errorq != nil {
		return "", errorq
	}

	request.Header.Set("Authorization", "Api-Key "+API_KEY)

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

	log.Println("Body:", objBody)
	log.Println(resp)

	return objBody.Id, nil
}

// Функция распознования речи
func recognize(filePath string, lang string) (string, error) {
	var recognitionText string = "Test: " + filePath + "\nLang: " + lang
	stringFileContent, err := uploadToBucket(filePath)

	if err != nil {
		log.Println("Error while uploading file to bucket:", err)
		return "", err
	}

	id, err := sendRequest(stringFileContent, lang)

	if err != nil {
		log.Println("Error while sending data: ", err)
		return "", err
	}

	recognitionText, err = getResult(id)

	if err != nil {
		log.Println("Error while getting result: ", err)
		return "", err
	}

	return recognitionText, nil
}
