package client

type Response struct {
	NormText string `json:"normText"`
	RawText  string `json:"rawText"`
	Error    string `json:"error,omitempty"`
	Details  string `json:"details,omitempty"`
}

type Error struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
