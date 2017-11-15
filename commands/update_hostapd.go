package commands

import (
	"os"
	"path/filepath"
)

type UpdateHostAPDCommand struct {
	destinationConfigFile string
	templateLocation      string
	hostAPDUnitName       string
}

type UpdateHostAPDRequest struct {
	NetworkSSID   string `json:"ssid"`
	Channel       int16  `json:"channel"`
	MACFiltering  bool   `json:"filter_mac"`
	BroadcastSSID bool   `json:"broadcast_ssid"`
	WPAPassphrase string `json:"wpa_passphrase"`
}

func (comm *UpdateHostAPDCommand) Init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	comm.templateLocation = filepath.Dir(ex) + "/example/hostapd.conf"
	comm.destinationConfigFile = "/tmp/hostapd.conf"
	comm.hostAPDUnitName = "hostapd.service"
}

//
func (comm *UpdateHostAPDCommand) Identifier() string {
	return "update-hostapd"
}

//
func (comm *UpdateHostAPDCommand) Object() interface{} {
	return &UpdateHostAPDRequest{}
}

//
func (comm *UpdateHostAPDCommand) Handle(command interface{}) []byte {
	req := command.(*UpdateHostAPDRequest)

	// Update Configuration File
	if err := handleTemplate(req, comm.templateLocation, comm.destinationConfigFile); err != nil {
		return hasError("Unable to update configuration file")
	}

	// Restart via SystemD
	if err := reloadOrRestartSystemdUnit(comm.hostAPDUnitName); err != nil {
		panic(err)
		return hasError("Unable to restart HostAPD service")
	}

	returnData := make(map[string]interface{})
	returnData["file_updated"] = comm.destinationConfigFile
	returnData["new_values"] = req

	return ranSuccessfully(returnData)
}

//
func (comm *UpdateHostAPDCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Update HostAPD"
	desc.Description = "Updates WiFi network broadcast settings, including SSID and Authentication."
	desc.Command = "update-hostapd"

	return *desc
}
