package normalization

type Request struct {
	Text   string `json:"text"`
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}
