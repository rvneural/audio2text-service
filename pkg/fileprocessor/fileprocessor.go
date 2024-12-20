package fileprocessor

import (
	"os"
	"path/filepath"
	"strings"

	"math/rand"

	"github.com/rs/zerolog"
)

type FileConverter interface {
	ConvertFile(filePath string, fileType string) (string, error)
}

type FileProcessor struct {
	Logger        *zerolog.Logger
	FileConverter FileConverter
}

func New(fileConvertor FileConverter, logger *zerolog.Logger) *FileProcessor {
	return &FileProcessor{Logger: logger, FileConverter: fileConvertor}
}

func (fp *FileProcessor) ProcessFile(fileData []byte, fileType string) (string, error) {
	fp.Logger.Info().Msg("FileProcessor: Processing file")
	filePath, err := fp.saveFile(fileData, fileType)
	if err != nil {
		return "", err
	}

	return fp.FileConverter.ConvertFile(filePath, fileType)
}

func (fp *FileProcessor) saveFile(fileData []byte, fileType string) (string, error) {
	fp.Logger.Info().Msg("Saving file")
	var filePath = "./../../uploads/" + fp.getRandonName(50) + "." + strings.ToLower(fileType)
	file, err := os.Create(filePath)

	if err != nil {
		fp.Logger.Error().Msg("Error creating file: " + err.Error())
		if os.IsNotExist(err) {
			err = os.MkdirAll("./../..uploads", 0755)
			if err != nil {
				fp.Logger.Error().Msg("Error creating directory: " + err.Error())
				return "", err
			}
			return fp.saveFile(fileData, fileType)
		}
		fp.Logger.Error().Msg("Error creating file: " + err.Error())
		return "", err
	}
	defer file.Close()
	_, err = file.Write(fileData)

	if err != nil {
		fp.Logger.Error().Msg("Error writing to file: " + err.Error())
		return "", err
	}

	return filepath.Abs(filePath)
}

func (fp *FileProcessor) getRandonName(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	rune_name := make([]rune, length)
	for i := range rune_name {
		rune_name[i] = letters[rand.Intn(len(letters))]
	}
	return string(rune_name)
}
