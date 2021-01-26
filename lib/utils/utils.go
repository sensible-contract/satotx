package utils

func ReverseBytes(data []byte) (result []byte) {
	for _, b := range data {
		result = append([]byte{b}, result...)
	}
	return result
}
