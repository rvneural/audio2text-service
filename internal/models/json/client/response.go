package client

type Response struct {
	NormText string `json:"normText"`
	RawText  string `json:"rawText"`
}

type Error struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
