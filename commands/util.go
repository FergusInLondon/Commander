// util.go - Fergus In London <fergus@fergus.london> (November 2017)
package commands

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	"github.com/godbus/dbus"
)

// A Command is an object that is capable of handling an instruction from the
//  JSON API.
type Command interface {
	// Init allows the Command to perform initialisation of dependencies
	//  and handle
	Init()

	// Identifier returns the short alphanumeric string which allows a client
	//  to call the action.
	Identifier() string

	// Description provides basic information on a given command, such as
	//  the name, the command identifier, and human readable description.
	Description() CommandDescription

	// DisplaySchema returns a JSON Schema describing the configuration
	//  object. Useful for validation and form building.
	DisplaySchema(w http.ResponseWriter, req *http.Request)

	// RetrieveConfig retrieves the current configuration parameters,
	//  allowing the client to display a form, or perform checks.
	RetrieveConfig(w http.ResponseWriter, req *http.Request)

	// SubmitConfig handles POST actions to the Command, generally for
	//  the submission of new configuration options.
	SubmitConfig(w http.ResponseWriter, req *http.Request)
}


// CommandDescription holds basic information about a command which is supplied
//  via the HTTP JSON API when a client hits the "/listing" endpoint.
type CommandDescription struct {
	// Name is a human readable name for the command - i.e "DNS Configuration"
	Name        string `json:"name"`

	// Command is the alphanumeric string used to call the command - i.e "dnsconf"
	Command     string `json:"command"`

	// Description is a human readable description, generally one or two sentences.
	Description string `json:"description"`
}

// Helper functions for returning standard payloads describing success and errors.

func hasError(writer http.ResponseWriter, msg string) {
	payload := make(map[string]interface{})
	payload["success"] = false
	payload["error"] = msg

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : false }")
	}

	writer.WriteHeader(500)
	writer.Write(jsonObject)
}

func ranSuccessfully(writer http.ResponseWriter, result interface{}) {
	payload := make(map[string]interface{})
	payload["success"] = true
	payload["result"] = result

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : true }")
	}

	writer.WriteHeader(200)
	writer.Write(jsonObject)
}


// Helper functions for interfacing with underlying configuration systems
//  such as DBus or Configuration Files.

func handleTemplate(data interface{}, tplFile, destination string) (err error) {
	tmpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return
	}

	file, err := os.Create(destination)
	if err != nil {
		return
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}


var dbusConnections struct {
	System  *dbus.Conn
	Session *dbus.Conn
}


func getDbusSystemConnection() (conn *dbus.Conn, err error) {
	if dbusConnections.System == nil {
		dbusConnections.System, err = dbus.SystemBus()
	}

	return dbusConnections.System, err
}


func getDbusSessionConnection() (conn *dbus.Conn, err error) {
	if dbusConnections.Session == nil {
		dbusConnections.Session, err = dbus.SessionBus()
	}

	return dbusConnections.Session, err
}

func reloadOrRestartSystemdUnit(unitName string) (err error) {
	dbusConn, err := getDbusSystemConnection()
	if err != nil {
		panic(err)
	}

	systemdManager := dbusConn.Object("org.freedesktop.systemd1", dbus.ObjectPath("/org/freedesktop/systemd1"))
	call := systemdManager.Call("org.freedesktop.systemd1.Manager.ReloadOrRestartUnit", 0, unitName, "replace")

	return call.Err
}
