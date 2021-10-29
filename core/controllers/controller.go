package controllers

import "github.com/obonobo/mux"

type Controller interface {
	Attach(router *mux.Router)
}
