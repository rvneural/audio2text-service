package main

func recognize(filePath string, lang string) string {
	var recognitionText string = "Test: " + filePath + "\nLang: " + lang

	return recognitionText
}
