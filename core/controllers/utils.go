package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(status int, rw http.ResponseWriter, data map[string]interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	rw.Write(append(j, '\n'))
	return nil
}

func MustWriteJSON(status int, rw http.ResponseWriter, data map[string]interface{}) {
	err := WriteJSON(status, rw, data)
	if err != nil {
		panic(err)
	}
}

func WriteOkJSON(rw http.ResponseWriter, data map[string]interface{}) error {
	return WriteJSON(http.StatusOK, rw, data)
}

func MustWriteOkJSON(rw http.ResponseWriter, data map[string]interface{}) {
	MustWriteJSON(http.StatusOK, rw, data)
}

func CreateSubrouter(router *mux.Router, prefix string) *mux.Router {
	if prefix == "" || prefix == "/" {
		return router
	}
	return router.PathPrefix(prefix).Subrouter()
}

func OptionsFunc(methods string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// rw.Write()
	}
}
