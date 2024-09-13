package post

type Request struct {
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

		TextNormalization struct {
			TextNormalization string `json:"textNormalization"`
		} `json:"textNormalization"`

		AudioProcessingType string `json:"audioProcessingType"`
	} `json:"recognitionModel"`

	SpeakerLabeling struct {
		SpeakerLabeling string `json:"speakerLabeling"`
	} `json:"speakerLabeling"`
	URI string `json:"uri"`
}
