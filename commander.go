// commander.go - Fergus In London <fergus@fergus.london> (November 2017)
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FergusInLondon/Commander/commands"
)

// Commander is responsible for the parsing of Command requests via HTTP,
//  and the subsequent dispatch to the correct reciever (or Command object).
type Commander struct {
	initialised bool
	registry    map[string]commands.Command
	status      struct {
		Uptime    time.Time `json:"initialisation_time"`
		Successes int64     `json:"successful_executions"`
		Failures  int64     `json:"failed_executions"`
	}
}

type commandHolder struct {
	Command    string          `json:"command"`
	Parameters json.RawMessage `json:"parameters"`
}

// Register places a Command in to the registry
func (comm *Commander) Register(c commands.Command) {
	comm.registry[c.Identifier()] = c
}

// Init configures the command registry, and begins the uptime counter.
func (comm *Commander) Init() (err error) {
	comm.registry = make(map[string]commands.Command)

	// @todo - find a nicer way.
	commObjects := make([]commands.Command, 4)
	commObjects[0] = new(commands.EchoCommand)
	commObjects[1] = new(commands.NotifyCommand)
	commObjects[2] = new(commands.ServicesCommand)
	commObjects[3] = new(commands.UpdateHostAPDCommand)

	for i := 0; i < len(commObjects); i++ {
		commObjects[i].Init()
		comm.Register(commObjects[i])
	}

	comm.status.Uptime = time.Now()
	comm.initialised = true
	return
}

// HandleCommand recieves a HTTP Request containing a JSON object with a "command"
//  string, and a "parameters object". It uses the "command property" to determine
//  which Command object is required, then unmarshels the "parameters" object in
//  to the correct form of struct.
func (comm *Commander) HandleCommand(w http.ResponseWriter, req *http.Request) {
	var payload commandHolder

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)
	defer req.Body.Close()

	if err != nil {
		comm.status.Failures++
		panic("Unable to parse command object.")
	} else {
		comm.status.Successes++
	}

	if commandHandler, ok := comm.registry[payload.Command]; ok {
		var commandParams = commandHandler.Object()
		json.Unmarshal(payload.Parameters, commandParams)

		response := commandHandler.Handle(commandParams)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

// Listing returns a JSON payload containing the number of registered commands,
//  as well as their names, a description, and the command string.
func (comm *Commander) Listing(w http.ResponseWriter, req *http.Request) {
	payload := make(map[string]interface{})
	descriptions := []commands.CommandDescription{}

	for _, command := range comm.registry {
		descriptions = append(descriptions, command.Description())
	}

	payload["count"] = len(descriptions)
	payload["commands"] = descriptions

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : false ")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObject)
}

// ProvideStatus responds to HTTP requests and returns a small JSON payload
//  detailing a few metrics about Commander's status; namely: (a) the uptime,
//  (b) number of commands handled, and (c) number of successful and failed commands.
func (comm *Commander) ProvideStatus(w http.ResponseWriter, req *http.Request) {
	status := make(map[string]interface{})
	status["uptime"] = time.Since(comm.status.Uptime).String()

	executions := make(map[string]int64)
	executions["successful"] = comm.status.Successes
	executions["failed"] = comm.status.Failures
	executions["total"] = (comm.status.Successes + comm.status.Failures)
	status["executions"] = executions

	jsonObject, err := json.Marshal(status)
	if err != nil {
		comm.status.Failures++
		jsonObject = []byte("{ \"success\" : false ")
	} else {
		comm.status.Successes++
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonObject)
}
