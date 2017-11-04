// notify.go - Fergus In London <fergus@fergus.london> (November 2017)
//
// A slightly more reasonable example - this command accepts a JSON payload
//  with a "title" and "message" property, and uses it to invoke a notification
//  via the users Desktop Environment.
//
// This is more as a proof of concept that Commander could interface with network
//  interfaces and other lower level APIs via DBus; I'm aware that the CoreOS
//  systemd lib I'm using also has dbus functionality though: so I may rework this.
package commands

import "github.com/godbus/dbus"

type ServicesCommand struct {
	systemdManager dbus.BusObject
}

type ServicesRequest struct {
}

func (sc *ServicesCommand) Init() {
	dbusConnection, err := getDbusSystemConnection()
	if err != nil {
		panic(err)
	}

	sc.systemdManager = dbusConnection.Object("org.freedesktop.systemd1", dbus.ObjectPath("/org/freedesktop/systemd1"))
}

//
func (sc *ServicesCommand) Identifier() string {
	return "services"
}

//
func (sc *ServicesCommand) Object() interface{} {
	return make(map[string]interface{})
}

//
func (sc *ServicesCommand) Handle(command interface{}) []byte {
	//req := command.(map)

	call := sc.systemdManager.Call(
		// Object Method
		"org.freedesktop.systemd1.Manager.ListUnits",
		// Flags
		0)

	if call.Err != nil {
		panic(call.Err)
		return hasError("Unable to query systemd")
	}

	returnData := make(map[string]interface{})
	returnData["services"] = call.Body
	returnData["notified_via"] = "dbus"
	return ranSuccessfully(returnData)
}

//
func (sc *ServicesCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Services"
	desc.Description = "Returns a list of running services via systemd."
	desc.Command = "services"

	return *desc
}
