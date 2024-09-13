package converter

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func (fc *FileConverter) convertVideo(filePath string, fileType string) (string, error) {
	var audioPath string = strings.Replace(filePath, "."+fileType, ".wav", -1)

	// ffmpeg -i video.flv audio.mp3
	cmd := exec.Command("ffmpeg", "-i", filePath, audioPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	// Возвращаем ошибку, если произошел сбой при конвертации
	if err != nil {
		fc.Logger.Error().Msg("Error converting video to audio: " + out.String())
		return audioPath, err
	}

	// Пытаемся удалить исходный файл, чтобы не задваивать занимаемое пространство
	go os.Remove(filePath)

	fc.Logger.Info().Msg("Converted video to: " + audioPath)

	return audioPath, nil
}
