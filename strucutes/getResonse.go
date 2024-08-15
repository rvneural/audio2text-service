package strucutes

type StartEndTime struct {
	StartTimeMs string `json:"startTimeMs"`
	EndTimeMs   string `json:"endTimeMs"`
}

type Words struct {
	Text string `json:"text"`
	StartEndTime
}

type Alternatives struct {
	Words []Words `json:"words"`
	Text  string  `json:"text"`
	StartEndTime
	Confidence float64     `json:"confidence"`
	Languages  []Languages `json:"languages"`
}

type Languages struct {
	LanguageCode string  `json:"languageCode"`
	Probabiliti  float64 `json:"probabiliti"`
}

type Labels struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

type Quantiles struct {
	Level float64 `json:"level"`
	Value float64 `json:"value"`
}

type PerSecond struct {
	Min       float64     `json:"min"`
	Max       float64     `json:"max"`
	Mean      float64     `json:"mean"`
	Std       float64     `json:"std"`
	Quantiles []Quantiles `json:"quantiles"`
}

type SpeakerInterrupts struct {
	SpeakerTag           string         `json:"speakerTag"`
	InterruptsCount      string         `json:"interruptsCount"`
	InterruptsDurationMs string         `json:"interruptsDurationms"`
	Interrupts           []StartEndTime `json:"interrupts"`
}

type Result struct {
	SessionUuid struct {
		Uuid          string `json:"uuid"`
		UserRequestId string `json:"userRequestId"`
	} `json:"sessionUuid"`
	AudioCursors struct {
		ReceivedDataMs string `json:"receivedDataMs"`
		ResetTimeMs    string `json:"resetTimeMs"`
		PartialTimeSm  string `json:"partialTimeSm"`
		FinalTimeMs    string `json:"finalTimeMs"`
		FinalIndex     string `json:"finalIndex"`
		EouTimeMs      string `json:"eouTimeMs"`
	} `json:"audiocursors"`
	ResponseWallTimeMs string `json:"responseWallTimeMs"`
	ChannelTag         string `json:"channelTag"`
	Partial            struct {
		Alternatives []Alternatives `json:"alternatives"`
		ChannelTag   string         `json:"channelTag"`
	} `json:"partial"`
	Final struct {
		Alternatives []Alternatives `json:"alternatives"`
		ChannelTag   string         `json:"channelTag"`
	} `json:"final"`
	EouUpdate struct {
		TimeMs string `json:"timeMs"`
	} `json:"eouUpdate"`
	FinalRefinement struct {
		FinalIndex     string `json:"finalIndex"`
		NormalizedText struct {
			Alternatives []Alternatives `json:"alternatives"`
			ChannelTag   string         `json:"channelTag"`
		} `json:"normalizedText"`
	} `json:"finalRefinement"`
	StatusCode struct {
		CodeType string `json:"codeType"`
		Message  string `json:"message"`
	} `json:"statusCode"`
	ClassifierUpdate struct {
		WindowType       string `json:"windowType"`
		StartTimeMs      string `json:"startTimeMs"`
		EndTimeMs        string `json:"endTimeMs"`
		ClassifierResult struct {
			Classifier string   `json:"classifier"`
			Highlights []Words  `json:"highlights"`
			Labels     []Labels `json:"labels"`
		} `json:"classifierResult"`
	} `json:"classifierUpdate"`
	SpeakerAnalysis struct {
		SpeakerTag       string `json:"speakerTag"`
		WindowType       string `json:"windowType"`
		SpeechBoundaries struct {
			StartEndTime
		} `json:"cspeechBoudaries"`
		TotalSpeechMs  string  `json:"totalSpeechMs"`
		SpeechRatio    float64 `json:"speechRation"`
		TotalSilenceMs string  `json:"totalSilenceMs"`
		SilenceRatio   float64 `json:"silenceRatio"`
		WordsCount     string  `json:"wordsCount"`
		LettersCount   string  `json:"lettersCount"`
		WordsPerSecond struct {
			PerSecond
		} `json:"wordsPerSecond"`
		LettersPerSecond struct {
			PerSecond
		} `json:"lettersPerSecond"`
		WordsPerUtterance struct {
			PerSecond
		} `json:"wordsPerUtterance"`
		LettersPerUtterance struct {
			PerSecond
		} `json:"lettersPerUtterance"`
		UtteranceCount              string `json:"utteranceCount"`
		UtteranceDurationEstimation struct {
			PerSecond
		} `json:"utteranceDurationEstimation"`
	} `json:"speakerAnalysis"`
	ConversationAnalysis struct {
		ConversationBoundaries struct {
			StartEndTime
		} `json:"conversationBoudaries"`
		TotalSimultaneousSilenceDurationMs    string  `json:"totalSimultaneousSilenceDurationMs"`
		TotalSimultaneousSilenceRatio         float64 `json:"totalSimultaneousSilenceRatio"`
		SimultaneousSilenceDurationEstimation struct {
			PerSecond
		} `json:"simultaneousSilenceDurationEstimation"`
		TotalSimultaneousSpeechDurationMs    string  `json:"totalSimultaneousSpeechDurationMs"`
		TotalSimultaneousSpeechRatio         float64 `json:"totalSimultaneousSpeechRatio"`
		SimultaneousSpeechDurationEstimation struct {
			PerSecond
		} `json:"simultaneousSpeechDurationEstimation"`
		SpeakerInterrupts     []SpeakerInterrupts `json:"speakerInterrupts"`
		TotalSpeechDurationMs string              `json:"totalSpeechDurationMs"`
		TotalSpeechRatio      float64             `json:"totalSpeechRatio"`
	} `json:"conversationAnalysis"`
}

type GetResponse struct {
	Result Result `json:"result"`
}
