package utils

import (
	"bytes"
	"encoding/json"

	"github.com/gorilla/mux"
)

func CreateSubrouter(router *mux.Router, prefix string) *mux.Router {
	if prefix == "" || prefix == "/" {
		return router
	}
	return router.PathPrefix(prefix).Subrouter()
}

// Create some compact json data. Calls `json.Compact`
func Compact(data string) []byte {
	var buf bytes.Buffer
	json.Compact(&buf, []byte(data))
	return buf.Bytes()
}
