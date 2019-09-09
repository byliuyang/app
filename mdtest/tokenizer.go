package mdtest

import (
	"encoding/json"

	"github.com/byliuyang/app/fw"
)

type CryptoTokenizer struct {
}

func (t CryptoTokenizer) Encode(payload fw.TokenPayload) (string, error) {
	buf, err := json.Marshal(payload)
	return string(buf), err
}

func (t CryptoTokenizer) Decode(tokenStr string) (fw.TokenPayload, error) {
	payload := map[string]interface{}{}
	err := json.Unmarshal([]byte(tokenStr), &payload)
	return payload, err
}

var CryptoTokenizerFake fw.CryptoTokenizer = CryptoTokenizer{}
