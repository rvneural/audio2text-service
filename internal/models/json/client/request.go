package client

type Request struct {
	Operation_ID string `json:"operation_id" xml:"operation_id" form:"operation_id"`
	URL          string `json:"url" xml:"url" form:"url"`
	Model        string `json:"model" xml:"model" form:"model"`
	File         struct {
		Name string `json:"name" xml:"name" form:"name"`
		Data []byte `json:"data" xml:"data" form:"data"`
		Type string `json:"type" xml:"type" form:"type"`
	} `json:"file" xml:"file" form:"file"`
	Languages []string `json:"languages" xml:"languages" form:"languages"`
	Dialog    bool     `json:"dialog" xml:"dialog" form:"dialog"`
	UserID    string   `json:"user_id" xml:"user_id" form:"user_id"`
}

type DBResult struct {
	FileName string `json:"filename"`
	RawText  string `json:"raw_text"`
	NormText string `json:"norm_text"`
}
