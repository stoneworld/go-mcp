package pkg

import (
	"fmt"

	"github.com/bytedance/sonic"
)

var sonicAPI = sonic.Config{UseInt64: false}.Froze() // 可有效防止整型溢出

func JsonUnmarshal(data []byte, v interface{}) error {
	if err := sonicAPI.Unmarshal(data, v); err != nil {
		return fmt.Errorf("%w: data=%s, err=%w", ErrJsonUnmarshal, data, err)
	}
	return nil
}
