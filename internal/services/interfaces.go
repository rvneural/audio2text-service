package services

type Recognition interface {
	RecognizeAudio(filePath string, lang []string, dialog bool, uniqPhraseSplitter string, maxLength int) ([]string, error)
}

type Normalization interface {
	NormalizeText(string) string
}

type FileProcessor interface {
	ProcessFile(fileData []byte, fileType string) (string, error)
}
