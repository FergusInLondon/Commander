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

type NotifyCommand struct {
}

type NotifyRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (nc *NotifyCommand) Init() {

}

//
func (ec *NotifyCommand) Identifier() string {
	return "notify"
}

//
func (ec *NotifyCommand) Object() interface{} {
	return &NotifyRequest{}
}

//
func (ec *NotifyCommand) Handle(command interface{}) []byte {
	req := command.(*NotifyRequest)

	dbusConnection, err := dbus.SessionBus()
	if err != nil {
		return hasError("Unable to configure dbus connection.")
	}

	notifier := dbusConnection.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	status := notifier.Call(
		// Object Method
		"org.freedesktop.Notifications.Notify",
		// Flags
		0,
		// DBus Method Params - 1: App Name
		"",
		// DBus Method Params - 2: Replaces Id
		uint32(0),
		// DBus Method Params - 3: App Icon
		"",
		// DBus Method Params - 4: Summary
		req.Title,
		// DBus Method Params - 4: Body
		req.Message,
		// DBus Method Params - 6: Actions
		[]string{},
		// DBus Method Params - 7: Hints
		map[string]dbus.Variant{},
		// DBus Method Params - 8: Expiry/Timeout
		int32(5000))

	if status.Err != nil {
		return hasError("Unable to create notification")
	}

	returnData := make(map[string]interface{})
	returnData["title"] = req.Title
	returnData["message"] = req.Message
	returnData["notified_via"] = "dbus"
	return ranSuccessfully(returnData)
}

//
func (ec *NotifyCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Notify"
	desc.Description = "Displays a system notification via the Desktop Environment."
	desc.Command = "notify"

	return *desc
}
