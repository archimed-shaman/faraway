package protocol

type ChallengeReq struct {
	Challenge  []byte `json:"challenge"`
	Difficulty int    `json:"difficulty"`
}

type ChallengeResp struct {
	Challenge  []byte `json:"challenge"`
	Difficulty int    `json:"difficulty"`
	Solution   []byte `json:"solution"`
}

type ErrorResp struct {
	Reason string `json:"reason"`
}

type Data struct {
	Payload []byte `json:"payload"`
}
