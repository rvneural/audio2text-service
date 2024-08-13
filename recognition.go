package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// Функция разделения аудиофайла на фрагменты определенной длины
func splitFile(filePath string, dir string, maxDuraion int) error {
	// ffmpeg -i somefile.mp3 -f segment -segment_time 3 -c copy out%03d.mp3
	cmd := exec.Command("ffmpeg", "-i", filePath, "-f", "segment", "-segment_time", strconv.FormatInt(int64(maxDuraion), 10), "-c", "copy", dir+"out%03d.wav")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Получаем массивы байт в видео строки из содержимого файлов
func readFileContent(pathToPempFiles string) ([]string, error) {
	var stringFileContent []string

	// Получаем содержимое папки
	entries, err := os.ReadDir(pathToPempFiles)
	if err != nil {
		log.Println("Error during reading files:", err)
		return stringFileContent, err
	}

	// Читаем каждый файл в массив байт
	byteFileContent := make([][]byte, len(entries))
	for i, e := range entries {
		file, err := os.ReadFile(pathToPempFiles + e.Name())
		if err != nil {
			log.Println("Errur during reading file:", err)
			return stringFileContent, err
		}
		byteFileContent[i] = file
	}

	// Переносим байты в строки as-is
	var wg sync.WaitGroup
	stringFileContent = make([]string, len(byteFileContent))
	for i := range len(byteFileContent) {
		wg.Add(1)
		go func(s *string) {
			defer wg.Done()
			var buffer bytes.Buffer
			for b, q := range byteFileContent[i] {
				buffer.WriteString(strconv.Itoa(int(q)))
				if b != len(byteFileContent[i])-1 {
					buffer.WriteString(" ")
				}
			}
			*s = buffer.String()
		}(&stringFileContent[i])
	}
	wg.Wait()

	return stringFileContent, nil
}

// Функция распознования речи
func recognize(filePath string, lang string) string {
	var recognitionText string = "Test: " + filePath + "\nLang: " + lang
	var pathToPempFiles = getTemFileDir()
	var maxTime int = 25

	err := splitFile(filePath, pathToPempFiles, maxTime)

	if err != nil {
		os.RemoveAll(pathToPempFiles)
		return err.Error()
	}

	stringFileContent, err := readFileContent(pathToPempFiles)

	if err != nil {
		return err.Error()
	}

	//os.RemoveAll(pathToPempFiles)
	return recognitionText
}
