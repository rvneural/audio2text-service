package structures

type RequestFromMainServer struct {
	FilePath string `json:"filePath"`
	Dialog   bool   `json:"dialog"`
	Language string `json:"language"`
}
