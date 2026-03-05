package json

import "github.com/bytedance/sonic"

type JSON struct {
}

func NewJSON() *JSON {
	return &JSON{}
}

func (j *JSON) Unmarshal(buf []byte, val interface{}) error {
	return sonic.Unmarshal(buf, val)
}
