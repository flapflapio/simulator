package schemacontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flapflapio/simulator/core/controllers/utils"
	"github.com/flapflapio/simulator/core/simulation/automata"
	"github.com/flapflapio/simulator/core/simulation/machine"
	"github.com/obonobo/mux"
)

const schemaFilename = "machine.schema.json"

type SchemaController struct {
	prefix string
}

func New() *SchemaController {
	return &SchemaController{}
}

func WithPrefix(prefix string) *SchemaController {
	return &SchemaController{
		prefix: prefix,
	}
}

func (sc *SchemaController) Attach(router *mux.Router) {
	r := utils.CreateSubrouter(router, sc.prefix)
	r.Methods("GET").Path("/machine.schema.json").HandlerFunc(Schema)
	r.Methods("POST").Path("/validate").HandlerFunc(Validate)
}

// If successful: 200 + machine json.
// If the machine in request body is invalid: 422.
func Validate(rw http.ResponseWriter, r *http.Request) {
	m, err := automata.Load(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write([]byte(m.Json()))
	if err != nil {
		panic(err)
	}
}

func Schema(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Del("Content-Disposition")
	rw.Header().Add(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", schemaFilename))

	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")

	rw.WriteHeader(200)
	rw.Write(getSchema())
}

func getSchema() []byte {
	schema, err := machine.GetSchema()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	return data
}
