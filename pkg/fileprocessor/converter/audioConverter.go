package converter

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func (fc *FileConverter) convertAudio(filePath string, fileType string) (string, error) {
	var audioPath string

	// Если текущий файл с беспроблемным типом, то конвертируем его в моно
	if fileType == "wav" {
		audioPath = strings.Replace(filePath, "."+strings.ToLower(fileType), "-mono.wav", -1)
	} else {
		audioPath = strings.Replace(filePath, "."+strings.ToLower(fileType), ".wav", -1)
	}

	// Создаем и выполняем cmd команду с ffmpeg, которая конвертирует аудио в .mp3
	cmd := exec.Command("ffmpeg", "-i", filePath, "-ac", "1", audioPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	// Возвращаем ошибку, если произошел сбой при конвертации
	if err != nil {
		fc.Logger.Error().Msg("Error converting audio to wav: " + out.String())
		return audioPath, errors.New("ffmpeg: " + err.Error())
	}

	// Пытаемся удалить исходный файл, чтобы не задваивать занимаемое пространство
	go os.Remove(filePath)

	fc.Logger.Info().Msg("Converted audio to: " + audioPath)

	return audioPath, nil
}
