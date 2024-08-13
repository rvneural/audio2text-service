package main

import "os"

func recognize(filePath string, lang string) string {
	var recognitionText string = "Test: " + filePath + "\nLang: " + lang
	var pathToPempFiles = getTemFileDir()

	os.RemoveAll(pathToPempFiles)
	return recognitionText
}
