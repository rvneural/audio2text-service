package structures

type SendRequest struct {
	RecognitionModel struct {
		Model       string
		AudioFormat struct {
			ContainerAudio struct {
				ContainerAudioType string `json:"containerAudioType"`
			} `json:"containerAudio"`
		} `json:"audioFormat"`
		LanguageRestriction struct {
			RestrictionType string   `json:"restrictionType"`
			LanguageCode    []string `json:"languageCode"`
		} `json:"languageRestriction"`
	} `json:"recognitionModel"`
	URI string `json:"uri"`
}
