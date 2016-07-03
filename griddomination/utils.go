package griddomination

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"github.com/pquerna/ffjson/ffjson"
	"errors"
)

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	buf, _ := ffjson.Marshal(&data)
	w.Write(buf)
	ffjson.Pool(buf)
}

func locationFromId(id string) (int64, int64, error) {
	coords := strings.Split(id, ".")

	if len(coords) != 2 {
		return 0, 0, errors.New("invalid format")
	}

	x, err := strconv.ParseInt(coords[0], 10, 64)

	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.ParseInt(coords[1], 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}

// credit: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	letterIdxBits = 6                    // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1 << letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func generateSessionToken() string {
	n := 32
	result := make([]byte, n)
	bufferSize := int(float64(n) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < n; j++ {
		if j % bufferSize == 0 {
			randomBytes = SecureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j % n] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result)
}

func SecureRandomBytes(n int) []byte {
	var randomBytes = make([]byte, n)
	rand.Read(randomBytes)
	return randomBytes
}
