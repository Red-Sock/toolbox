package toolbox

import (
	"crypto/rand"
	"encoding/base64"
)

func Random(length int) ([]byte, error) {
	randKey := make([]byte, length)
	_, err := rand.Read(randKey)
	if err != nil {
		return nil, err
	}

	key := make([]byte, base64.StdEncoding.EncodedLen(len(randKey)))

	base64.StdEncoding.Encode(key, randKey)
	return key[:length], nil
}
