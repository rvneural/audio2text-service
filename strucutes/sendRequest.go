package strucutes

type SendRequest struct {
	RecognitionModel struct {
		Model       string
		AudioFormat struct {
			ContainerAudio struct {
				ContainerAudioType string `json:"containerAudioType"`
			} `json:"containerAudio"`
		} `json:"audioFormat"`
		TextNormalization struct {
			TextNormalization   string `json:"textNormalization"`
			ProfanityFilter     bool   `json:"profanityFilter"`
			LiteratureText      bool   `json:"literatureText"`
			PhoneFormattingMode string `json:"phoneFormattingMode"`
		} `json:"textNormalization"`
		LanguageRestriction struct {
			RestrictionType string   `json:"restrictionType"`
			LanguageCode    []string `json:"languageCode"`
		} `json:"languageRestriction"`
		AudioProcessingType string `json:"audioProcessingType"`
	} `json:"recognitionModel"`
	RecognitionClassifier struct {
		Classifiers []Classifiers `json:"classifiers"`
	} `json:"recognitionClassifier"`
	SpeechAnalysis struct {
		EnableSpeakerAnalysis          bool     `json:"enableSpeakerAnalysis"`
		EnableConversationAnalysis     bool     `json:"enableConversationAnalysis"`
		DescriptiveStatisticsQuantiles []string `json:"descriptiveStatisticsQuantiles"`
	} `json:"speechAnalysis"`
	SpeakerLabeling struct {
		SpeakerLabeling string `json:"speakerLabeling"`
	} `json:"speakerLabeling"`
	Content string `json:"uri"`
}
