package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
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
