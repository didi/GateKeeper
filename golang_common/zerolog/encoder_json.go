// +build json_log

package zerolog

// encoder_json.go file contains bindings to generate
// JSON encoded byte stream.

import (
	"github.com/didi/gatekeeper/golang_common/zerolog/internal/json"
)

var (
	_ encoder = (*json.Encoder)(nil)

	enc = json.Encoder{}
)

func appendJSON(dst []byte, j []byte) []byte {
	return append(dst, j...)
}

func decodeIfBinaryToString(in []byte) string {
	return string(in)
}

func decodeObjectToStr(in []byte) string {
	return string(in)
}

func decodeIfBinaryToBytes(in []byte) []byte {
	return in
}
