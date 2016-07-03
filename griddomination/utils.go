package griddomination

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
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

func generateSessionToken() string {
	n := 64
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(b)
}

type Location struct {
	X int64
	Y int64
}

func LocationFromId(id string) *Location {
	coords := strings.Split(id, ".")

	if len(coords) != 2 {
		return nil
	}

	x, err := strconv.ParseInt(coords[0], 10, 64)

	if err != nil {
		return nil
	}

	y, err := strconv.ParseInt(coords[1], 10, 64)

	if err != nil {
		return nil
	}

	return &Location{X:x, Y:y}
}

func (location *Location) ToId() string {
	return fmt.Sprintf("%v.%v", location.X, location.Y)
}
