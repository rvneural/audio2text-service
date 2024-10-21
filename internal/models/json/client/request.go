package client

type Request struct {
	URL  string `json:"url" xml:"url" form:"url"`
	File struct {
		Data []byte `json:"data" xml:"data" form:"data"`
		Type string `json:"type" xml:"type" form:"type"`
	} `json:"file" xml:"file" form:"file"`
	Languages []string `json:"languages" xml:"languages" form:"languages"`
	Dialog    bool     `json:"dialog" xml:"dialog" form:"dialog"`
}
