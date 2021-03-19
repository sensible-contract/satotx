package utils

func ReverseBytes(data []byte) (result []byte) {
	dataLen := len(data)
	result = make([]byte, dataLen)
	for idx, b := range data {
		result[dataLen-1-idx] = b
	}
	return result
}
