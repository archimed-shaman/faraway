package protocol

import "encoding/json"

type NonceReq struct{}

type NonceResp struct {
	Nonce      []byte `json:"nonce"`
	Difficulty int    `json:"difficulty"`
}

type DataReq struct {
	Nonce      []byte `json:"nonce"`
	Difficulty int    `json:"difficulty"`
	CNonce     []byte `json:"cnonce"`
}

type DataResp struct {
	Payload []byte `json:"payload"`
}

type ErrorResp struct {
	Reason string `json:"reason"`
}

type Package struct {
	Tag     string          `json:"tag"`
	Payload json.RawMessage `json:"payload"`
}
