package simulationcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/flapflapio/simulator/core/app"
	"github.com/flapflapio/simulator/core/controllers/utils"
	"github.com/flapflapio/simulator/core/simulation"
	"github.com/flapflapio/simulator/core/simulation/automata"
	"github.com/gorilla/websocket"
	"github.com/obonobo/mux"
)

const (
	INVALID_MACHINE_MSG = `` +
		`{"Err":"The machine that was sent is not ` +
		`valid or otherwise could not be processed"}`

	PLEASE_PROVIDE_A_TAPE_MSG = `` +
		`{"Err":"Please provide an input with query param 'tape'"}`

	FAILED_TO_CREATE_A_NEW_SIMULATION = `` +
		`{"Err":"Failed to create a new simulation"}`

	FAILED_TO_OBTAIN_RESULTS_OF_SIMULATION = `` +
		`{"Err":"Failed to obtain results of simulation"}`

	FAILED_TO_CREATE_A_RESPONSE = `` +
		`{"Err":"Failed to create a response"}`
)

type SimulationController struct {
	prefix    string
	simulator simulation.Simulator
}

func New(simulator simulation.Simulator) *SimulationController {
	return &SimulationController{
		prefix:    "/",
		simulator: simulator,
	}
}

// Attaches this controller to the given router
func (c *SimulationController) Attach(router *mux.Router) {
	r := utils.CreateSubrouter(router, c.prefix)
	r.Methods("POST").Path("/simulate").HandlerFunc(c.DoSimulation)
	r.Methods("DELETE").Path("/simulation/{id}").HandlerFunc(c.EndSimulation)
	r.Methods("POST").Path("/simulation/start").HandlerFunc(c.StartSimulation)
	r.Methods("GET").Path("/ws").HandlerFunc(c.WebSocket)
}

func (c *SimulationController) StartSimulation(rw http.ResponseWriter, r *http.Request) {
	m, err := automata.Load(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(INVALID_MACHINE_MSG))
		return
	}

	r.ParseForm()
	tape, ok := r.Form["tape"]
	if !ok || len(tape) < 1 {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(PLEASE_PROVIDE_A_TAPE_MSG))
		return
	}

	id, err := c.simulator.Start(m, tape[0])
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(FAILED_TO_CREATE_A_NEW_SIMULATION))
		return
	}

	rw.Header().Add("content-type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusAccepted)
	rw.Write([]byte(fmt.Sprintf(`{"Status":"Accepted","Id":%v}`, id)))
}

func (c *SimulationController) WithPrefix(prefix string) *SimulationController {
	return &SimulationController{
		prefix:    app.Trim(prefix),
		simulator: c.simulator,
	}
}

func (c *SimulationController) EndSimulation(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(`{"Err":"Simulation id cannot be empty"}`))
		return
	}

	intVar, err := strconv.Atoi(id)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(`{"Err":"Simulation id is not valid"}`))
		return
	}

	err = c.simulator.End(intVar)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(fmt.Sprintf(`{"Err":"%v"}`, err)))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"Status":"Simulation ended successfully"}`))
}

func (c *SimulationController) DoSimulation(rw http.ResponseWriter, r *http.Request) {
	m, err := automata.Load(r.Body)

	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(INVALID_MACHINE_MSG))
		log.Println(err)
		return
	}

	r.ParseForm()
	tape, ok := r.Form["tape"]
	if !ok || len(tape) < 1 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(PLEASE_PROVIDE_A_TAPE_MSG))
		return
	}

	var sim simulation.Simulation

	// Create a new simulation
	id, err := c.simulator.Start(m, tape[0])
	if check(err, rw, FAILED_TO_CREATE_A_NEW_SIMULATION) {
		return
	}

	// Run the simulation from start to finish
	for sim = c.simulator.Get(id); !sim.Done(); sim.Step() {
	}

	// Grab the result of the simulation
	res, err := sim.Result()
	if check(err, rw, FAILED_TO_OBTAIN_RESULTS_OF_SIMULATION) {
		return
	}

	// Serialize result
	data, err := json.Marshal(res)
	if check(err, rw, FAILED_TO_CREATE_A_RESPONSE) {
		return
	}

	// Write result to response body
	rw.Header().Del("Content-Type")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(append(data, '\n'))
	c.simulator.End(id)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketMessage struct {
	Op      string                 `json:"op"`
	Params  map[string]string      `json:"params"`
	Payload map[string]interface{} `json:"payload"`
}

// This endpoint receives websocket connections
func (c *SimulationController) WebSocket(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return // Upgrade will reply automatically with an error response
	}
	defer conn.Close()

	// For now, the server simply waits for requests on the connection and
	// responds with the result of the operation requested

	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			unexpected := websocket.IsUnexpectedCloseError(
				err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			if unexpected {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		// Trimming the message a little bit
		message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))

		var wsMessage WebSocketMessage
		err = json.Unmarshal(message, &wsMessage)

		checkWsErr := func(err error) bool {
			if err != nil {
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Err: %v", err)))
				return true
			}
			return false
		}

		if err != nil {
			err = fmt.Errorf("invalid command")
		}

		if checkWsErr(err) {
			continue
		}

		// Otherwise log the message
		log.Println(wsMessage)

		switch wsMessage.Op {
		case "DO":
			m, err := automata.Load(wsMessage.Payload)
			if checkWsErr(err) {
				continue
			}

			tape, ok := wsMessage.Params["tape"]
			if !ok || len(tape) < 1 {
				if checkWsErr(fmt.Errorf("please provide an input tape in the 'params' JSON field of your request")) {
					continue
				}
			}

			var sim simulation.Simulation

			// Create a new simulation
			id, err := c.simulator.Start(m, tape)
			if checkWsErr(err) {
				continue
			}
			defer c.simulator.End(id)

			// Run the simulation from start to finish
			for sim = c.simulator.Get(id); !sim.Done(); sim.Step() {
			}

			// Grab the result of the simulation
			res, err := sim.Result()
			if checkWsErr(err) {
				continue
			}

			// Serialize result
			data, err := json.Marshal(res)
			if checkWsErr(err) {
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			conn.WriteMessage(websocket.TextMessage, data)
		case "START":
			m, err := automata.Load(wsMessage.Payload)
			if checkWsErr(err) {
				continue
			}

			tape, ok := wsMessage.Params["tape"]
			if !ok || len(tape) < 1 {
				if checkWsErr(fmt.Errorf("please provide an input tape in the 'params' JSON field of your request")) {
					continue
				}
			}

			id, err := c.simulator.Start(m, tape)
			if checkWsErr(err) {
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"Status":"Accepted","Id":%v}`, id)))
		case "END":
			id, ok := wsMessage.Params["id"]
			if !ok || len(id) < 1 {
				if checkWsErr(fmt.Errorf("please provide a simulation id in the 'params' JSON field of your command")) {
					continue
				}
			}

			intVar, err := strconv.Atoi(id)
			if err != nil {
				checkWsErr(fmt.Errorf("id is not valid"))
			}

			err = c.simulator.End(intVar)
			if err != nil {
				checkWsErr(fmt.Errorf("id is not valid"))
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			conn.WriteMessage(websocket.TextMessage, []byte(`{"Status":"Simulation ended successfully"}`))
		default:
			checkWsErr(fmt.Errorf("unsupported operation '%v'", wsMessage.Op))
		}
	}
}

func check(err error, rw http.ResponseWriter, msg string) bool {
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(msg))
		return true
	}
	return false
}
