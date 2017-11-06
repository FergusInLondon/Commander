package commands

import (
	"os"
	"path/filepath"
)

type UpdateDNSMasqCommand struct {
	destinationConfigFile string
	templateLocation      string
	DNSMasqUnitName       string
}

type UpdateDNSMasqRequest struct {
	DHCPRangeBegin string   `json:"dhcp_begin"`
	DHCPRangeEnd   string   `json:"dhcp_end"`
	DHCPLeaseTime  string   `json:"dhcp_lease"`
	DNSServers     []string `json:"dhcp_servers"`
}

func (comm *UpdateDNSMasqCommand) Init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	comm.templateLocation = filepath.Dir(ex) + "/example/dnsmasq.conf"
	comm.destinationConfigFile = "/tmp/dnsmasq.conf"
	comm.DNSMasqUnitName = "uuidd.service"
}

//
func (comm *UpdateDNSMasqCommand) Identifier() string {
	return "update-dnsmasq"
}

//
func (comm *UpdateDNSMasqCommand) Object() interface{} {
	return &UpdateDNSMasqRequest{}
}

//
func (comm *UpdateDNSMasqCommand) Handle(command interface{}) []byte {
	req := command.(*UpdateDNSMasqRequest)

	// Update Configuration File
	if err := handleTemplate(req, comm.templateLocation, comm.destinationConfigFile); err != nil {
		panic(err)
		return hasError("Unable to update configuration file")
	}

	// Restart via SystemD
	if err := reloadOrRestartSystemdUnit(comm.DNSMasqUnitName); err != nil {
		panic(err)
		return hasError("Unable to restart DNSMasq service")
	}

	returnData := make(map[string]interface{})
	returnData["file_updated"] = comm.destinationConfigFile
	returnData["new_values"] = req

	return ranSuccessfully(returnData)
}

//
func (comm *UpdateDNSMasqCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Update dnsmasq"
	desc.Description = "Updates the DNS and DHCP settings, controlled via dnsmasq."
	desc.Command = "update-dnsmasq"

	return *desc
}
