// commander.go - Fergus In London <fergus@fergus.london> (November 2017)
package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/fergusinlondon/Commander/commands"
	)

// Commander is responsible for the parsing of Command requests via HTTP,
//  and the subsequent dispatch to the correct reciever (or Command object).
type Commander struct {
	registry    []commands.Command
	status      struct {
		Uptime    time.Time `json:"initialisation_time"`
		Successes int64     `json:"successful_executions"`
		Failures  int64     `json:"failed_executions"`
	}
}


// Init configures the command registry, and begins the uptime counter.
func (comm *Commander) Init() (err error) {
	comm.registry = append(comm.registry,
		new(commands.EchoCommand),
		new(commands.NotifyCommand),
		new(commands.TemplateCommand),
		new(commands.ServicesCommand))

	for _, command := range comm.registry {
		command.Init()
	}

	comm.status.Uptime = time.Now()
	return
}


// Listing returns a JSON payload containing the number of registered commands,
//  as well as their names, a description, and the command string.
func (comm *Commander) Listing(w http.ResponseWriter, req *http.Request) {
	payload := make(map[string]interface{})
	descriptions := []commands.CommandDescription{}

	for _, command := range comm.registry {
		descriptions = append(descriptions, command.Description())
	}

	payload["commands"] = descriptions
	payload["count"] = len(descriptions)

	w.Header().Set("Content-Type", "application/json")
	if jsonObject, err := json.Marshal(payload); err != nil {
		w.Write([]byte("{ \"success\" : false "))
	} else {
		w.Write(jsonObject)
	}
}


// ProvideStatus responds to HTTP requests and returns a small JSON payload
//  detailing a few metrics about Commander's status; namely: (a) the uptime,
//  (b) number of commands handled, and (c) number of successful and failed commands.
func (comm *Commander) ProvideStatus(w http.ResponseWriter, req *http.Request) {
	executions := make(map[string]int64)
	executions["failed"] = comm.status.Failures
	executions["successful"] = comm.status.Successes
	executions["total"] = (comm.status.Successes + comm.status.Failures)

	status := make(map[string]interface{})
	status["executions"] = executions
	status["registed_commands"] = len(comm.registry)
	status["uptime"] = time.Since(comm.status.Uptime).String()

	w.Header().Set("Content-Type", "application/json")
	if jsonObject, err := json.Marshal(status); err != nil {
		comm.status.Failures++
		w.Write([]byte("{ \"success\" : false "))
	} else {
		comm.status.Successes++
		w.Write(jsonObject)
	}
}


// ConfigureListeners configures the mux router for responding to inbound
//  API requests.
func (comm *Commander) ConfigureListeners(router *mux.Router)  {
	router.HandleFunc("/listing", comm.Listing)
	router.HandleFunc("/status", comm.ProvideStatus)

	actions := router.PathPrefix("/action/").
		Headers("X-Requested-With", "PiRouterBackend").Subrouter()

	for _, command := range comm.registry {
		commandAction := actions.PathPrefix(command.Identifier()).Subrouter()
		commandAction.HandleFunc("/describe", command.DisplaySchema)
		commandAction.Methods("GET").Path("/").HandlerFunc(command.RetrieveConfig)
		commandAction.Methods("POST").Path("/").HandlerFunc(command.SubmitConfig)
	}
}