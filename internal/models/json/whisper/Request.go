package whisper

type Request struct {
	FileData []byte `json:"file"`
	FileName string `json:"name"`
}
