package client

type Request struct {
	URL       string   `json:"url" xml:"url" form:"url"`
	FileData  []byte   `json:"fileData" xml:"fileData" form:"fileData"`
	Languages []string `json:"languages" xml:"languages" form:"languages"`
	FileType  string   `json:"fileType" xml:"fileType" form:"fileType"`
	Dialog    bool     `json:"dialog" xml:"dialog" form:"dialog"`
}
