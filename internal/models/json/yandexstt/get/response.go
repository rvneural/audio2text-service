package get

type Alternatives struct {
	Text string `json:"text"`
}

type Result struct {
	ChannelTag string `json:"channelTag"`
	Final      struct {
		Alternatives []Alternatives `json:"alternatives"`
	} `json:"final"`
}

type Response struct {
	Result Result `json:"result"`
}
