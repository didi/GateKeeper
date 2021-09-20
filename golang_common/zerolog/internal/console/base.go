package console

type Encoder struct{}

// AppendKey appends a new key to the output JSON.
func (e Encoder) AppendKey(dst []byte, key string) []byte {

	if len(dst) > 1 {
		dst = append(dst, []byte("||")...)
	}
	if key != "" {
		dst = e.AppendString(dst, key+"=")
	}

	return dst
}
