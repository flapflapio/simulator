package types

import "github.com/gorilla/mux"

type Controller interface {
	Attach(router *mux.Router)
}