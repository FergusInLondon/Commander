// util.go - Fergus In London <fergus@fergus.london> (November 2017)
package commands

import (
	"encoding/json"
	"os"
	"text/template"

	"github.com/godbus/dbus"
)

// A Command is an object that is capable of handling an instruction from the
//  JSON API.
type Command interface {
	Init()
	Identifier() string
	Object() interface{}
	Handle(interface{}) []byte
	Description() CommandDescription
}

// CommandDescription holds basic information about a command which is supplied
//  via the HTTP JSON API when a client hits the "/listing" endpoint.
type CommandDescription struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

func hasError(msg string) (jsonObject []byte) {
	payload := make(map[string]interface{})
	payload["success"] = false
	payload["error"] = msg

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : false }")
	}

	return
}

func ranSuccessfully(result map[string]interface{}) (jsonObject []byte) {
	payload := make(map[string]interface{})
	payload["success"] = true
	payload["result"] = result

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : true }")
	}

	return
}

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
