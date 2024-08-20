package structures

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
	Final              struct {
		Alternatives []Alternatives `json:"alternatives"`
		ChannelTag   string         `json:"channelTag"`
	} `json:"final"`
	StatusCode struct {
		CodeType string `json:"codeType"`
		Message  string `json:"message"`
	} `json:"statusCode"`
}

type GetResponse struct {
	Result Result `json:"result"`
}
