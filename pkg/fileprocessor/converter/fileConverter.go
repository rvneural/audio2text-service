package converter

import (
	"errors"
	"slices"
	"strings"

	"github.com/rs/zerolog"
)

type FileConverter struct {
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *FileConverter {
	return &FileConverter{Logger: logger}
}

func (fc *FileConverter) ConvertFile(filePath string, fileType string) (string, error) {
	// Доступные форматы видео
	var video_types = []string{"mp4", "mov", "avi", "mpg", "flv", "webm", "mkv", "m4v", "hevc"}

	// Доступные форматы аудио
	var audio_types = []string{"mp3", "m4a", "wav", "aac", "ogg", "flac", "opus", "m4r", "aiff", "wma"}

	fc.Logger.Info().Msg("Converting file: " + filePath)

	if slices.Contains(video_types, strings.ToLower(fileType)) {
		fc.Logger.Info().Msg("Starting convertion of video: " + filePath)
		return fc.convertVideo(filePath, fileType)
	} else if slices.Contains(audio_types, strings.ToLower(fileType)) {
		fc.Logger.Info().Msg("Starting convertion of audio: " + filePath)
		return fc.convertAudio(filePath, fileType)
	} else {
		return filePath, errors.New("incorrect file type")
	}
}
