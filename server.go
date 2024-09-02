package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rvRecognitionService/structures"
)

const (
	SERVER = "127.0.0.1"
	PORT   = "45679"
)

func mainRecognizer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println("Request from", r.RemoteAddr)
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error while reading body:", err)
		var ans structures.Resonse
		ans.NormText = err.Error()
		ans.RawText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	var request structures.RequestFromMainServer
	err = json.Unmarshal(data, &request)
	if err != nil {
		log.Println("Error unmarshalling body:", err)
		var ans structures.Resonse
		ans.NormText = err.Error()
		ans.RawText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	log.Println(r.RemoteAddr, " >>>", request)

	filePath := request.FilePath
	isDialog := request.Dialog
	language := request.Language

	log.Println("Расшифровка файла")
	rawText, normText, err := recognize(filePath, language, isDialog)
	log.Println(r.RemoteAddr, "->\tRaw text:", rawText)

	if err != nil {
		log.Println("Error while recognition:", err)
		var ans structures.Resonse
		ans.NormText = err.Error()
		ans.RawText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	var response structures.Resonse
	response.NormText = normText
	response.RawText = rawText

	byteResponse, err := json.Marshal(response)

	if err != nil {
		log.Println("Error while marshalling response:", err)
		var ans structures.Resonse
		ans.NormText = err.Error()
		ans.RawText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	w.WriteHeader(200)
	w.Write(byteResponse)
}

func startServer() {
	http.HandleFunc("/", mainRecognizer)
	log.Printf("Server started at http://%s:%s/\n", SERVER, PORT)

	// Start the main server
	err := http.ListenAndServe(SERVER+":"+PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
