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
	"rvRecognitionService/strucutes"
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
		log.Println("ERROR WHILE UPLOADING FILE TO BUCKET")
		log.Println(out.String())
		return "", err
	}

	bucketAddress = "https://storage.yandexcloud.net/" + BUCKET_NAME + "/" + key
	return bucketAddress, nil
}

func getResult(id string) (string, error) {
	httpURL := "http://stt.api.cloud.yandex.net/stt/v3/getRecognition?operationId=" + id

	request, errorq := http.NewRequest("GET", httpURL, nil)
	if errorq != nil {
		return "", errorq
	}
	//var getResponse strucutes.GetResponse

	request.Header.Set("Authorization", "Api-Key "+API_KEY)

	client := &http.Client{}

	var resp *http.Response

	var GetResponse strucutes.GetResponse
	for {
		log.Println("Asking for trasncription result")
		resp, _ = client.Do(request)
		if resp.StatusCode != 200 {
			if resp.StatusCode == 500 {
				data, _ := io.ReadAll(resp.Body)
				return "", errors.New("yandex can not recognize this file\n\nlog: " + string(data))
			}
			if resp.StatusCode == 400 {
				log.Println("Code:", resp.StatusCode)
				data, _ := io.ReadAll(resp.Body)
				return "", errors.New(string(data))
			}
			log.Println("Code:", resp.StatusCode)
			<-time.After(10 * time.Second)
			continue
		} else {
			log.Println("Code:", resp.StatusCode)
			break
		}
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var err error
	line, err := reader.ReadBytes('\n')

	if err != nil {
		log.Println("Error while reading result:", err)
		return "", err
	}

	err = json.Unmarshal(line, &GetResponse)
	if err != nil {
		log.Println("Error while unmarshalling getting response:", err)
		log.Println("\n\n\n>>>>>", string(line))
		return "", err
	}
	log.Println(GetResponse)

	return GetResponse.Result.Final.Alternatives[0].Text, nil
}

func sendRequest(fileData string, lang string) (string, error) {
	httpURL := "https://stt.api.cloud.yandex.net/stt/v3/recognizeFileAsync"
	httpMethod := "POST"

	// Инициализация тела запроса
	httpBody := strucutes.SendRequest{}
	httpBody.RecognitionModel.Model = "general:rc"

	httpBody.RecognitionModel.AudioFormat.ContainerAudio.ContainerAudioType = "WAV"
	httpBody.RecognitionModel.TextNormalization.TextNormalization = "TEXT_NORMALIZATION_DISABLED"
	httpBody.RecognitionModel.TextNormalization.ProfanityFilter = false
	httpBody.RecognitionModel.TextNormalization.LiteratureText = true
	httpBody.RecognitionModel.TextNormalization.PhoneFormattingMode = "PHONE_FORMATTING_MODE_DISABLED"
	httpBody.RecognitionModel.LanguageRestriction.RestrictionType = "WHITELIST"
	httpBody.RecognitionModel.LanguageRestriction.LanguageCode = []string{lang}
	httpBody.RecognitionModel.AudioProcessingType = "FULL_DATA"

	classifiere := strucutes.Classifiers{
		//Classifier: "formal_greeting",
		//Triggers:   []string{"ON_FINAL"},
		Classifier: "",
		Triggers:   []string{},
	}

	httpBody.RecognitionClassifier.Classifiers = []strucutes.Classifiers{classifiere}

	httpBody.SpeechAnalysis.EnableSpeakerAnalysis = true
	httpBody.SpeechAnalysis.EnableConversationAnalysis = true
	httpBody.SpeechAnalysis.DescriptiveStatisticsQuantiles = []string{}

	httpBody.SpeakerLabeling.SpeakerLabeling = "SPEAKER_LABELING_DISABLED"
	httpBody.Content = fileData

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
	log.Println("File sent")

	if err2 != nil {
		log.Println("Error wending file to server")
		return "", err2
	}
	defer resp.Body.Close()

	byteResultBodyF := make([]byte, 5000)

	log.Println("Reading sending file response")
	q, err := resp.Body.Read(byteResultBodyF)

	byteResultBody := byteResultBodyF[:q]

	if err != nil {
		log.Println("Error while reading sending response")
		return "", err
	}

	var objBody strucutes.SendResponse

	//// NEED TO COPY
	log.Println("Unmarshalling sending file response")
	err = json.Unmarshal(byteResultBody, &objBody)
	log.Println(objBody)

	if err != nil {
		log.Println("Error while unmarshalling sending response:", err)
		log.Println(string(byteResultBody))
		return "", err
	}

	return objBody.Id, nil
}

// Функция распознования речи
func recognize(filePath string, lang string) string {
	var recognitionText string = "Test: " + filePath + "\nLang: " + lang

	//stringFileContent, err := readFileContent(filePath)
	//if err != nil {
	//	log.Println("Error while getting string content:", err)
	//	return err.Error()
	//}

	stringFileContent, err := uploadToBucket(filePath)

	if err != nil {
		log.Println("Error while uploading file to bucket:", err)
		return err.Error()
	}

	id, err := sendRequest(stringFileContent, lang)

	if err != nil {
		log.Println("Error while sending data: ", err)
		return err.Error()
	}

	recognitionText, err = getResult(id)

	if err != nil {
		log.Println("Error while getting result: ", err)
		return err.Error()
	}

	//os.RemoveAll(pathToPempFiles)
	return recognitionText
}
