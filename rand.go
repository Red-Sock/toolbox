package toolbox

import (
	"encoding/base64"
	"math/rand"
)

var randKeys = map[byte]struct{}{}

func init() {
	baseStr := `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/`
	for _, k := range []byte(baseStr) {
		randKeys[k] = struct{}{}
	}
}

func RandomBase64(length int) []byte {
	randKey := make([]byte, 0, length)
	iL := length
	for iL > 0 {
		for k := range randKeys {
			randKey = append(randKey, k)
			iL--
			if rand.Intn(10) > 5 {
				break
			}
		}
	}

	key := make([]byte, base64.StdEncoding.EncodedLen(len(randKey)))

	base64.StdEncoding.Encode(key, randKey)
	return key[:length]
}
