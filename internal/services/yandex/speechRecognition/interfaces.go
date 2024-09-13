package speechRecognition

type Uploader interface {
	Upload(path string) (string, error)
}

type Recognizer interface {
	SendRequest(path string, lang []string, dialog bool) (string, error)
	GetResponse(id string) ([][]byte, error)
}

type Parser interface {
	Parse(rawData [][]byte, uniqPhraseSplitter string, maxLength int) ([]string, error)
}
