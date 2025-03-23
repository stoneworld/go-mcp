package pkg

import "github.com/bytedance/sonic"

var sonicAPI = sonic.Config{UseInt64: false}.Froze() // 可有效防止整型溢出

func JsonUnmarshal(data []byte, v interface{}) error {
	return sonicAPI.Unmarshal(data, v)
}
